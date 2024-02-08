package types

import "github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"

type (
	MetaData interface {
		GetKey() string
		GetValue() string
	}

	medaData struct {
		key   string
		value string
	}

	MetaDataDecorator struct {
		status *user.Status
	}
)

func (m *medaData) GetKey() string {
	return m.key
}

func (m *medaData) GetValue() string {
	return m.value
}

func (d *MetaDataDecorator) AddMetaData(meta MetaData) {
	d.status.Metadata[meta.GetKey()] = meta.GetValue()
}

func AppendMD(key, value string) MetaData {
	return &medaData{
		key:   key,
		value: value,
	}
}

func SetMetaData(status *user.Status, maps ...MetaData) {
	decorator := &MetaDataDecorator{
		status: status,
	}

	for _, data := range maps {
		decorator.AddMetaData(data)
	}
}
