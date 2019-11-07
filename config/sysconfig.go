package config

import (
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var Sysconfig = &sysConfig{}

func init() {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("Sys config read err")
	}
	err = jsoniter.Unmarshal(b, Sysconfig)
	if err != nil {
		panic(err)
	}
	go mkdir()

}
func mkdir() {
	_, err := os.Stat(Sysconfig.Basic.StaticPath)
	if err != nil { //不存在创建
		err = os.MkdirAll(Sysconfig.Basic.StaticPath, os.ModePerm)
		if err != nil {
			log.Println("创建文件夹失败")
			panic(err)
		}
	}
	_, err = os.Stat(Sysconfig.Basic.ReqLogPath)
	if err != nil { //不存在创建
		err = os.MkdirAll(Sysconfig.Basic.StaticPath, os.ModePerm)
		if err != nil {
			log.Println("创建文件夹失败")
			panic(err)
		}
	}
	_, err = os.Stat(Sysconfig.Basic.ErrorLogPath)
	if err != nil { //不存在创建
		err = os.MkdirAll(Sysconfig.Basic.StaticPath, os.ModePerm)
		if err != nil {
			log.Println("创建文件夹失败")
			panic(err)
		}
	}
}

type sysConfig struct {
	Basic    basic //基础信息
	Redis    redis
	Mysql    mysql
	RabbitMQ rabbitMQ
}
type basic struct {
	Port             string `json:"Port"`
	StaticPath       string `json:"StaticPath"`
	CertFile         string `json:"CertFile"`
	CertKey          string `json:"CertKey"`
	ReqLogPath       string
	ErrorLogPath     string
	IsLocal          bool //是否是本地 本地为true打印日志
	SocketCheckToken bool //websocket是否检查token
}

type redis struct {
	Password string
	Addr     string
}
type mysql struct {
	Name     string
	Host     string
	Port     string
	Database string
	Password string
}
type rabbitMQ struct {
	UserName string
	Password string
	Addr     string
}
