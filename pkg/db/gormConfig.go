package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
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
