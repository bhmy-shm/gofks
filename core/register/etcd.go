package register

import (
	"context"
	"fmt"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/errorx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"math/rand"
	"time"
)

// EtcdRegistry is etcd registry.
type EtcdRegistry struct {
	opts   *optionEtcd
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// NewEtcdRegistry creates etcd registry
func NewEtcdRegistry(config *gofkConfs.RegistryConfig, opts ...OptionEtcd) (r *EtcdRegistry) {

	client, err := NewEtcdClient(config)
	if err != nil {
		errorx.Fatal(err, "etcd连接失败")
		return nil
	}

	op := &optionEtcd{
		ctx:         context.Background(),
		namespace:   config.Namespace(),
		ttl:         config.TTL(),
		maxRetry:    config.MaxRetry(),
		dialTimeout: config.DialTimeout(),
	}
	for _, o := range opts {
		o(op)
	}
	return &EtcdRegistry{
		opts:   op,
		client: client,
		kv:     clientv3.NewKV(client),
	}
}

// NewEtcdClient 获取etcd Client
func NewEtcdClient(config *gofkConfs.RegistryConfig) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints(),
		DialTimeout: time.Second * config.DialTimeout(),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
}

func NewEtcdRegistryByClient(config *gofkConfs.RegistryConfig, client *clientv3.Client) (r *EtcdRegistry) {

	op := &optionEtcd{
		ctx:         context.Background(),
		namespace:   config.Namespace(),
		ttl:         config.TTL(),
		maxRetry:    config.MaxRetry(),
		dialTimeout: config.DialTimeout(),
	}

	return &EtcdRegistry{
		opts:   op,
		client: client,
		kv:     clientv3.NewKV(client),
	}
}

// Register the registration.
func (r *EtcdRegistry) Register(ctx context.Context, service *ServiceInstance) error {
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)
	value, err := marshal(service)
	if err != nil {
		return err
	}
	if r.lease != nil {
		r.lease.Close()
	}
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		return err
	}

	go r.heartBeat(r.opts.ctx, leaseID, key, value)
	return nil
}

// Deregister the registration.
func (r *EtcdRegistry) Deregister(ctx context.Context, service *ServiceInstance) error {
	defer func() {
		if r.lease != nil {
			r.lease.Close()
		}
	}()

	//todo 发送邮件告警消息
	//str := fmt.Sprintf(
	//	"Title: DataAccessIDS-【%s】\n dateTime: %s \n reason: %s",
	//	service.Name, time.Now().Format("2006-01-02 15:04"), "程序出现异常关闭，请检查",
	//)
	//err := mq.NewForeignMQ().SendMessage(Config.C.MQ2Exchange, []byte(str))
	//if err != nil {
	//	logger.Logf(logger.ERROR, fmt.Sprintf("mail send MQ failed= %v", err))
	//	log.Println("发送邮件告警消息失败")
	//	return err
	//}
	//log.Println("发送邮件告警消息成功")

	// 进行etcd 反向注册
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)
	_, err := r.client.Delete(ctx, key)
	return err
}

// GetService return the service instances in memory according to the service name.
func (r *EtcdRegistry) GetService(ctx context.Context, name string) ([]*ServiceInstance, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	resp, err := r.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]*ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		if si.Name != name {
			continue
		}
		items = append(items, si)
	}
	return items, nil
}

// Watch creates a watcher according to the service name.
func (r *EtcdRegistry) Watch(ctx context.Context, name string) (Watcher, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	return newWatcher(ctx, key, name, r.client)
}

// registerWithKV create a new lease, return current leaseID
func (r *EtcdRegistry) registerWithKV(ctx context.Context, key string, value string) (clientv3.LeaseID, error) {
	grant, err := r.lease.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return 0, err
	}
	_, err = r.client.Put(ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return 0, err
	}
	return grant.ID, nil
}

func (r *EtcdRegistry) heartBeat(ctx context.Context, leaseID clientv3.LeaseID, key string, value string) {
	curLeaseID := leaseID
	kac, err := r.client.KeepAlive(ctx, leaseID)
	if err != nil {
		curLeaseID = 0
	}
	rand.Seed(time.Now().Unix())

	for {
		if curLeaseID == 0 {
			// try to registerWithKV
			var retreat []int
			for retryCnt := 0; retryCnt < r.opts.maxRetry; retryCnt++ {
				if ctx.Err() != nil {
					return
				}
				// prevent infinite blocking
				idChan := make(chan clientv3.LeaseID, 1)
				errChan := make(chan error, 1)
				cancelCtx, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()
					id, registerErr := r.registerWithKV(cancelCtx, key, value)
					if registerErr != nil {
						errChan <- registerErr
					} else {
						idChan <- id
					}
				}()

				select {
				case <-time.After(3 * time.Second):
					cancel()
					continue
				case <-errChan:
					continue
				case curLeaseID = <-idChan:
				}

				kac, err = r.client.KeepAlive(ctx, curLeaseID)
				if err == nil {
					break
				}
				retreat = append(retreat, 1<<retryCnt)
				time.Sleep(time.Duration(retreat[rand.Intn(len(retreat))]) * time.Second)
			}
			if _, ok := <-kac; !ok {
				// retry failed
				return
			}
		}

		select {
		case _, ok := <-kac:
			if !ok {
				if ctx.Err() != nil {
					// channel closed due to context cancel
					return
				}
				// need to retry registration
				curLeaseID = 0
				continue
			}
		case <-r.opts.ctx.Done():
			return
		}
	}
}

func (r *EtcdRegistry) GetDialTimeout() time.Duration {
	return r.opts.dialTimeout
}
