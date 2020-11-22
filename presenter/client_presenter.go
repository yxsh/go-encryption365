package presenter

import (
	"errors"
	"fmt"
	"github.com/yxsh/go-encryption365/conf"
	"github.com/yxsh/go-encryption365/req"
	"os"
	"strconv"
	"strings"
)

type Presenter struct {
	client *req.Request365
	Config *conf.Config365
}
func New() *Presenter{
	p := new(Presenter)
	p.client = req.NewClient()
	return p
}
func (p *Presenter)SendCode(mail string) bool{
	result, err := p.client.SendAuthCode(mail)
	if err != nil {
		fmt.Println("发送错误:", err.Error())
		return false
	}
	if !result.IsOK() {
		fmt.Println("错误：", result.Message)
		return false
	}
	fmt.Println("已发送到:", mail, " 请检查邮箱，特别注意垃圾邮件处。")
	return true
}
func (p *Presenter)Register(info *req.RegisterInfo) bool{
	reg, err := p.client.AccountRegister(info)
	if err != nil {
		fmt.Println("注册发生了错误:", err.Error())
		return false
	}
	if !reg.IsOK() {
		fmt.Println("错误：", reg.Message)
		return false
	}
	return true
}
func (p *Presenter) RefreshToken(){
	if p.Config.AccessToken != "" {
		p.client.SetToken(p.Config.ClientId, p.Config.AccessToken)
	}
}
func (p *Presenter)Login(user,pwd string) bool{
	e := p.client.ClientCreate(user, pwd)
    if e != nil {
    	fmt.Println("登陆错误:", e)
    	return false
	}
    p.Config.UserName = user
    p.Config.AccessToken = p.client.Token.AccessToken
    p.Config.ClientId = p.client.Token.ClientId
    if e = p.Config.Save(); e != nil {
    	fmt.Println("登陆成功，但保存登陆Token失败了", e.Error())
	}
	return true
}
func (p *Presenter) GetOrderInfo(orderNum string) bool{
	order, e := p.client.GetOrderInfo(orderNum)
	if e != nil {
		fmt.Println(orderNum + "错误:", e.Error())
		return false
	}
	fmt.Println(order)
	return true
}
func (p *Presenter) GetOrder(orderNum string) (*req.OrderInfo365, error){
	return p.client.GetOrderInfo(orderNum)
}
func (p *Presenter) CertList() bool{
  ps,e := p.client.GetProducts()
  if e != nil {
  	fmt.Println("错误: ", e.Error())
  	return false
  }
  fmt.Println(ps)
  return true
}
func (p *Presenter) GenCSR(csr *req.CSR365) *req.CSR365{
 fmt.Println("正在生成，请勿操作")
 e := csr.Generate()
 if e != nil {
   fmt.Println("创建错误了", e.Error())
   return nil
 }
 return csr
}
func (p *Presenter) CreatNewCert(req *req.ReqCertInfo) (*req.OrderInfo365, bool){
	ret, e := p.client.CreateNewCert(req)
	if e != nil {
		fmt.Println("错误:", e.Error())
		return ret, false
	}
	id := strconv.Itoa(ret.TrustoceanId)
	e = p.Config.AddOrder(id)
	if e != nil && e != conf.ErrExistOrder {
		fmt.Println("添加订单记录失败,订单Id", ret.TrustoceanId)
		return ret, false
	}
	fmt.Println(ret)
	return ret, true
}
func (p *Presenter) SaveFile(name string, suffix string, val string) error{
	if name == "" {
		return errors.New("文件名为空")
	}
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "*", "")
	file, e := os.Create(name + "." + suffix)
	if e != nil {
		return e
	}
	defer file.Close()
	_, e = file.WriteString(val)
    if e != nil {
    	return e
	}
	return nil
}
func (p *Presenter) SaveCsrKey(name string, csr string, key string){
  p.SaveFile(name, "csr", csr)
  p.SaveFile(name, "key", key)
}
func (p *Presenter) AddOrderManager(orderNum string) bool{
	rs := p.GetOrderInfo(orderNum)
	if !rs {
		return false
	}
	e := p.Config.AddOrder(orderNum)
	if e != nil {
		fmt.Printf("订单[%s]添加错误:%s\r\n" , orderNum, e.Error())
		return false
	}
	return true
}