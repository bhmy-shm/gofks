package event

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/google/uuid"
	"log"
	"testing"
)

type testBody struct {
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
	Nid    string `json:"nid"`
}

type testParams struct {
	Id     int    `json:"id"`
	Method string `json:"method"`
	From   string `json:"from"`
	To     string `json:"to"`
	Params struct {
		Name   string `json:"name"`
		Mobile string `json:"mobile"`
		Email  string `json:"email"`
		Nid    string `json:"nid"`
	} `json:"params"`
}

func Test_message_Json_build(t *testing.T) {

	conf := &gofkConfs.MqConfig{}
	gofkConf.LoadConf(conf, gofkConf.WithPath("../../application.yaml"))

	tb := testBody{
		Name:   "shm",
		Mobile: "12312312312",
		Email:  "zdasd@qq.com",
		Nid:    uuid.New().String(),
	}
	message := NewMessage[*JsonMsg](
		AbleJsonMsg(),
		WithRequest(&JsonRequest{
			ID: 1,
			JsonBase: &JsonBase{
				Method: "contact.create",
				From:   "from-test",
				To:     "to-test",
			},
		}),
		WithParams(tb),
	)

	message.Pack() //通过NewMessage创建的对象，需要将设置的message进行打包，写入对应的序列化对象，jsonMsg

	log.Println("id:", message.Able.GetId())

	log.Println("message request:", *message.Able.GetRequest())

	log.Println("message params:", message.Able.GetParams())

	bb := message.Able.PackEncode() //这一步是序列化,jsonMsg 就是按照json进行序列化

	log.Println("encode end:", string(bb))

	// ------------------------ 将收到的消息解析成到指定对象中 -------------
	msg2 := testParams{}

	newMsg := NewMessage[*JsonMsg](AbleJsonMsg())

	err := newMsg.Able.PackDecode(bb, &msg2)
	if err != nil {
		log.Println(err)
	}
	log.Println("msg2 params Decode 2 interface:", msg2)

	// ---------------------- 再把解析之后的转换成 json ----------------
	//jsonMsg := newMsg.UnPackFromInterface(&msg2)
	//bytes := jsonMsg.Able.PackEncode()
	//
	//log.Println("msg2 Unpack and PackEncode：", string(bytes))
}

func Test_message_Json_UnPack(t *testing.T) {
	conf := &gofkConfs.MqConfig{}
	gofkConf.LoadConf(conf, gofkConf.WithPath("../../application.yaml"))

	tb := testBody{
		Name:   "shm",
		Mobile: "12312312312",
		Email:  "zdasd@qq.com",
		Nid:    uuid.New().String(),
	}
	message := NewMessage[*JsonMsg](
		AbleJsonMsg(),
		WithRequest(&JsonRequest{
			ID: 1,
			JsonBase: &JsonBase{
				Method: "contact.create",
				From:   "from-test",
				To:     "to-test",
			},
		}),
		WithParams(tb),
	)

	message.Pack()                  //通过NewMessage创建的对象，需要将设置的message进行打包，写入对应的序列化对象，jsonMsg
	bb := message.Able.PackEncode() //这一步是序列化,jsonMsg 就是按照json进行序列化

	// ------------

	msg := message.UnPackFromBytes(bb)
	log.Println("MSG:", string(msg.Able.GetMessage()))
}

func Test_message_Xml_build(t *testing.T) {

	//message := NewMessage[*XMLMsg](
	//	AbleXmlMsg(),
	//	WithId(2),
	//	WithOperator("increment"),
	//	WithType("notify"),
	//	WithBody(),
	//)
	//
	//message.Pack()
	//
	//print("mid:", message.Able.getId())
	//
	//for _, v := range message.Able.Body.Items {
	//	log.Println("item,v:", v)
	//}
	//
	//pack := message.Able.PackMsg()
	//log.Println("pack:", string(pack))
}
