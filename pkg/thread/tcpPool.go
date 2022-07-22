package thread

import (
	"errors"
	"fmt"
	"github.com/bhmy-shm/gofks/pkg/errorx"
	"sync"
	"time"
)

//线程池的配置信息

type Config struct {
	InitialCap  int                         //最小连接数
	MaxCap      int                         //最大并发存活连接数
	MaxIdle     int                         //最大空闲连接
	Factory     func() (interface{}, error) //生成连接
	Close       func(interface{}) error     //关闭连接
	Ping        func(interface{}) error     //检查连接是否有效
	IdleTimeout time.Duration               //连接最大的空闲时间，超过这个时间要丢弃连接
}

//线程池连接信息

type idleConn struct {
	conn interface{} //连接信息
	t    time.Time   //起始连接时间
}

type ChanPool struct {
	mx           sync.Mutex
	conns        chan *idleConn              //连接管道
	factory      func() (interface{}, error) //创建连接
	close        func(interface{}) error     //关闭连接
	ping         func(interface{}) error     //检查连接
	idleTimeout  time.Duration               //连接空闲超时时间
	waitTimeout  time.Duration               //等待时长
	maxActive    int                         //最大连接数
	openingConns int                         //已经连接的连接数
}

//可以借鉴sync.pool

type Pool interface {
	Get() (interface{}, error)
	Put(interface{}) error
	Close(interface{}) error
	Ping(interface{}) error
	Release()
	Len() int
}

func NewChannelPool(poolConfig *Config) (Pool, error) {
	if !(poolConfig.InitialCap <= poolConfig.MaxIdle &&
		poolConfig.MaxCap >= poolConfig.MaxIdle &&
		poolConfig.InitialCap >= 0) {
		return nil, errors.New("无效的容量设置")
	}
	if poolConfig.Factory == nil {
		return nil, errors.New("无效的功能设置，创建连接")
	}
	if poolConfig.Close == nil {
		return nil, errors.New("无效的功能设置，关闭连接")
	}

	//生成连接池对象
	c := &ChanPool{
		conns:        make(chan *idleConn, poolConfig.MaxIdle), //开辟全部空闲连接内存
		factory:      poolConfig.Factory,
		close:        poolConfig.Close,
		idleTimeout:  poolConfig.IdleTimeout,
		maxActive:    poolConfig.MaxCap,
		openingConns: poolConfig.InitialCap,
	}

	if poolConfig.Ping != nil {
		c.ping = poolConfig.Ping
	}

	//为最小连接数创建连接
	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, err := c.factory()
		if err != nil {
			return nil, fmt.Errorf("连接池未能正常填充最小连接数", err)
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}
	return c, nil
}

//获取连接池所有连接
func (c *ChanPool) getConns() chan *idleConn {
	c.mx.Lock()
	defer c.mx.Unlock()
	//声明一个临时变量返回
	conns := c.conns
	return conns
}

//检查连接是否正常
func (c *ChanPool) Ping(conn interface{}) error {
	if conn == nil {
		return errors.New("连接被拒绝")
	}
	return c.ping(conn)
}

//取出一个连接
func (c *ChanPool) Get() (interface{}, error) {
	//拿到连接池
	conns := c.getConns()
	if conns == nil {
		return nil, errorx.ErrClosed
	}

	//遍历连接池,如果存在连接，判断后直接返回该连接。如果没有连接，创建一个连接并返回。
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, errorx.ErrClosed
			}
			//判断当前连接是否超时，超时则直接丢弃连接(用临时变量进行判断)
			if timeout := c.idleTimeout; timeout > 0 {
				//如果起始连接时间 + 空闲超时时间 在 当前时间之前，代表已经超时
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					c.Close(wrapConn.conn)
					continue
				}
			}

			//判断当前连接是否失效
			if c.ping != nil {
				if err := c.Ping(wrapConn.conn); err != nil {
					c.close(wrapConn.conn)
					continue
				}
			}

			//如果都没有问题，则把这个连接返回
			return wrapConn.conn, nil

		default:
			//如果没有从连接池中取出数据
			c.mx.Lock()
			defer c.mx.Unlock()

			//如果当前连接数 >= 最大连接数，则代表连接池满了
			if c.openingConns >= c.maxActive {
				return nil, errorx.ErrMaxActiveConnReached
			}
			//如果还没有创建
			if c.factory == nil {
				return nil, errorx.ErrClosed
			}
			//直接创建一个连接
			conn, err := c.factory()
			if err != nil {
				return nil, err
			}
			c.openingConns++

			return conn, nil
		}
	}
}

//将连接重新入池
func (c *ChanPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is null don't in pool")
	}
	c.mx.Lock()

	//如果chan，连接池是空的,那还放个锤子了
	if c.conns == nil {
		c.mx.Unlock()
		return c.close(conn)
	}

	//重新生成连接
	select {
	case c.conns <- &idleConn{conn: conn, t: time.Now()}:
		c.mx.Unlock()
		return nil
	default:
		c.mx.Unlock()
		return c.Close(conn)
	}
}

//关闭某一个连接
func (c *ChanPool) Close(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is null don't in pool")
	}
	c.mx.Lock()
	defer c.mx.Unlock()

	if c.close == nil {
		return nil
	}
	c.openingConns--
	return c.close(conn)
}

func (c *ChanPool) Release() {
	c.mx.Lock()

	//将对象全部清空
	conns := c.conns
	closeFun := c.Close

	c.conns = nil
	c.factory = nil
	c.ping = nil
	c.close = nil
	c.mx.Unlock()

	if conns == nil {
		return
	}
	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn)
	}
}

//获取已有连接数量
func (c *ChanPool) Len() int {
	return len(c.getConns())
}
