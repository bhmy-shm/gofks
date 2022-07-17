package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type UserConfig map[interface{}]interface{}

//递归读取用户配置文件

func GetConfigValue(m UserConfig,prefix []string,index int) interface{}  {
	key:=prefix[index]
	if v,ok:=m[key];ok{
		if index==len(prefix)-1{ //到了最后一个
			return v
		}else{
			index=index+1
			if mv,ok:=v.(UserConfig);ok{ //值必须是UserConfig类型
				return GetConfigValue(mv,prefix,index)
			}else{
				return  nil
			}
		}
	}
	return  nil
}

type ServerConfig struct {
	Port int
}

type SysConfig struct {
	Server *ServerConfig
	Config UserConfig
}


func loadConfigFile() []byte {
	dir, _ := os.Getwd()
	file := dir + "/application.yaml"
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return nil
	}
	return b
}

func DefaultSysConfig() *SysConfig {
	return &SysConfig{
		Server: &ServerConfig{
			Port: 8080,
		},
	}
}


func InitSysConfig() *SysConfig {
	config := DefaultSysConfig()

	//读取配置文件
	if buf := loadConfigFile(); buf != nil {
		err := yaml.Unmarshal(buf, config)
		if err != nil {
			log.Fatal(err)
		}
	}
	return config
}
