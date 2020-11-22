package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	ConfigName  = "config365.json"
	FilePerm = 0755
)

/**
自动续签
 */
type Renew365 struct {
	Order string
	DestCert string
	DestKey string
}
/**
订单详情
 */
type Order365 struct {
	OceanId string
}
/**
基础配置
 */
type Config365 struct {
	UserName string `desc:"用户名" json:"user_name"`
	ClientId string `desc:"客户端ID" json:"client_id"`
	AccessToken string `desc:"APIToken" json:"access_token"`
	Orders []string `desc:"订单编号" json:"orders"`
}
var ErrExistOrder = errors.New("order Exists")
func Load() *Config365{
	config := new(Config365)
	info, e := os.Stat(ConfigName)
	if e != nil{
		return config
	}
	if info.IsDir() {
		fmt.Println("配置文件不应该为文件夹", info.Name())
		os.Exit(0)
	}
	file, e := os.Open(ConfigName)
	if e != nil {
		fmt.Println("读取配置文件是失败(权限?)", ConfigName)
		return config
	}
	defer file.Close()
	confText, e:= ioutil.ReadAll(file)
	e = json.Unmarshal(confText, config)
	if e != nil {
		fmt.Println("配置文件读取错误", e.Error())
		os.Exit(0)
	}
	return config
}
func (c *Config365)AddOrder(o string) error{
	if o == "" {
		return errors.New("order empty")
	}
	if c.Orders == nil || cap(c.Orders) < 1{
		c.Orders = make([]string, 1)
	}
	for _,d := range c.Orders {
		if d == o {
			return ErrExistOrder
		}
	}
	c.Orders = append(c.Orders, o)
	return c.Save()
}
func (c *Config365)Save() error{
	jsonFile,e := os.OpenFile(ConfigName,os.O_CREATE | os.O_WRONLY | os.O_TRUNC, FilePerm)
	if e != nil {
		return e
	}
	defer jsonFile.Close()
	data, e := json.Marshal(c)
	if e != nil {
		return e
	}
	return ioutil.WriteFile(ConfigName, data, FilePerm)
}

