package db

import (
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"reflect"
)

type GormAdapter struct {
	*gorm.DB
	//写入操作
	//查询操作
}

func GormPlug() *GormAdapter {

	dsn := "root:123456@tcp(127.0.0.1:3306)/corev2?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, //自动配置
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln(err)
	}

	//mysql连接池配置
	myDb, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	myDb.SetMaxIdleConns(5)
	myDb.SetMaxOpenConns(10)

	return &GormAdapter{DB: db}
}

func (this *GormAdapter) tableName() string {
	return "测试返回表名字"
}

func (this *GormAdapter) FindById(model interface{}, id int) []byte {
	newStruct := reflect.New(reflect.TypeOf(model))

	var (
		data = newStruct.Interface()
	)

	err := this.DB.Model(newStruct.Interface()).Where("id = ?", id).Find(&data).Error
	if err != nil {
		log.Fatalln("findById is err=", err)
	}

	buf, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("data marshal to string is failed", err)
	}
	return buf
}

func (this *GormAdapter) FindBySerial(model interface{}, serial string) []byte {
	return nil
}
