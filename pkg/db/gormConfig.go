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

type DbConfig struct{}

func NewDbConfig() *DbConfig {
	return &DbConfig{}
}

func (this *DbConfig) GormDb() *gorm.DB {
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
	return db
}

func (this *DbConfig) Injector() string {
	return "dbConfig"
}

func (this *DbConfig) FindById(model interface{}, id int) []byte {
	newStruct := reflect.New(reflect.TypeOf(model))

	var (
		data = newStruct.Interface()
	)

	err := this.GormDb().Model(newStruct.Interface()).Where("id = ?", id).Find(&data).Error
	if err != nil {
		log.Fatalln("findById is err=", err)
	}

	buf, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("data marshal to string is failed", err)
	}
	return buf
}
