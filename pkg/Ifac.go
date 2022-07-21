package pkg

//interface{} value 数据类型转换
type Value interface {
	Bool() (bool, error)
	Int() (int, error)
	String() (string, error)
	Float64() (float64, error)
	Json() (string, error)
	Interface() interface{}
	Slice() ([]interface{}, error)
	StringSlice() []string

	Str() string
}

//文件监听
type Watcher interface {
	Next() (*File, error)
	Stop() error
}

//json\yaml 编解码
type Encoder interface {
	Encode(interface{}) ([]byte, error) //编码
	Decode([]byte, interface{}) error   //解码
	String() string                     //返回文件名字
}

//读取文件的源内容
type Source interface {
	Read() (*File, error)
	Watch() (Watcher, error)
	String() string
	//Write(*ChangeSet) error
}
