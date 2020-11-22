package ui

import (
	"fmt"
	"github.com/yxsh/go-encryption365/presenter"
	"github.com/yxsh/go-encryption365/req"
	"os"
	"strconv"
)
var Presenter *presenter.Presenter

func submitOrder()  {
	fmt.Println("通过  https://console.trustocean.com/clientarea.php  登陆控制台")
	fmt.Println("选择左侧【管理证书】菜单，右侧表格显示的编号，即为订单编号")
	fmt.Println("本次只支持一个一个录入")
	orderNum := ComInput("请输入订单编号")
	if !Presenter.AddOrderManager(orderNum) {
	   if !ComSelect("是否新添加?") {
	   	  ShowMenu()
	   	  return
	   }
	}
	if ComSelect("继续添加?"){
		submitOrder()
	}else{
		ShowMenu()
	}
}
func showOrders(){
   orders := Presenter.Config.Orders
   if len(orders) < 1 {
   	  fmt.Println("无订单信息,如有历史订单，请使用 主菜单[5. 提交订单（由于没有获取所有订单接口，只能手工提交）] ")
   	  ShowMenu()
   	  return
   }
   for _,d := range orders {
   	 if d == "" {
   	 	continue
	 }
   	  Presenter.GetOrderInfo(d)
   }
   ShowMenu()
}
func showCertList(){
	Presenter.CertList()
	ShowMenu()
}
func genCSR(isShow bool) *req.CSR365{
	fmt.Print(`

   签证书的时候如果没有提供CSR，也会调用这里的方法生成
   !!注意：该功能调用 https://www.csr.sh 站点生成，如果不放心请自行提供即可!!

`)
	if !ComSelect("已阅读并知晓，继续？") {
		ShowMenu()
	}
	fmt.Print(`
----------------------------------------
     1.域名证书
     2.IP证书
     3.随便（随便填写证书信息）  
  1,2 两者区别在于 选择IP证书，会默认在域名列填写
  common-name-for-public-ip-address.com
  (签证书的时候，在填入IP即可)
----------------------------------------
`)
	sel := ComInputNum("请选择")

	if sel < 1 || sel > 3 {
		fmt.Println("错误的选择")
		return nil
	}
	var csr = req.NewCsr()
	if sel == 1 {
		csr = new(req.CSR365)
		InputPkg(csr)
	}
	if sel == 2 {
		csr = new(req.CSR365)
		csr.CommonName = "common-name-for-public-ip-address.com"
		InputPkg(csr)
	}
	ret :=  Presenter.GenCSR(csr)
	if csr != nil && isShow {
		fmt.Println("---------------CSR--------------")
		fmt.Println(ret.PublicCer)
		ComAnyKeyContinue()
		fmt.Println("---------------私钥--------------")
		fmt.Println(ret.PrivateKey)
	}
	return ret
}
func downloadCert(order *req.OrderInfo365){
   id := "0"
  if order == nil {
  	id = ComInput("输入订单Id")
  }else{
  	id = strconv.Itoa(order.TrustoceanId)
  }
  fmt.Println("正在获取订单编号[" ,id,"]","的证书信息...")
  order, e := Presenter.GetOrder(id)
  if e != nil {
	 fmt.Println("错误",id,":", e.Error())
	 if !ComSelect("是否重试?") {
	 	ShowMenu()
	 }else{
	 	downloadCert(order)
	 }
	 return
  }
  if order.CertStatus != "issued_active" || order.CertCode == "" {
  	if ComSelect("域名还未验证或者状态错误,是否重试") {
  		downloadCert(order)
	}else{
		ShowMenu()
	}
	return
  }
  for key,_ := range order.DcvIndex {
  	e := Presenter.SaveFile(key , "cer", order.CertCode)
  	if e != nil {
  		fmt.Println("下载证书错误,订单号",order.TrustoceanId,":", e.Error())
	}
  }
  fmt.Print(`
    ==========================================
      [成功] 已下载到当前目录，域名/IP.cer"
    =========================================
  `)
  ShowMenu()
}
func newCert(){
   fmt.Print(`
    -----------------签证书----------------
    谨慎操作：（有以下步骤）
     i   生成或者提供CSR
     ii  IP/域名 验证
     iii 下载证书(在本程序同一目录下) 
    --------------------------------------
`)
 csr := new(req.CSR365)
 if ComSelect("是否由本程序生成(还是粘贴已有的CSR)") {
	csr = genCSR(false)
	if csr == nil {
		ShowMenu()
	}
 }else{
 	fmt.Print("\r\n\r\n粘贴CSR,保留\r\n-----BEGIN CERTIFICATE REQUEST-----和-----END CERTIFICATE REQUEST-----\r\n:")
 	p := ComLongText("-----END CERTIFICATE REQUEST-----")
	fmt.Print("温馨提示: 请保管好私钥")
    csr.PublicCer = p
 }
 if !Presenter.CertList() {
	 ShowMenu()
 }
 sel := ComInput("请输入证书的产品编号(目前免费的为100)")
 ps := req.ProductIndex[sel]
 if ps == nil {
 	fmt.Println("编号",sel, "证书不存在！")
 	ShowMenu()
 	return
 }
 info := &req.ReqCertInfo{
 	 ProductId: ps.Id,
 	 Period: ps.Period,
 	 CsrCode: csr.PublicCer,
 }
 InputPkg(info)
 order, result := Presenter.CreatNewCert(info)
 if !result {
 	ShowMenu()
 	return
 }
for key,_ := range order.DcvIndex{
	Presenter.SaveCsrKey(key, csr.PublicCer, csr.PrivateKey)
}
 fmt.Println("如果已经设置好验证信息，继续...")
 ComAnyKeyContinue()
 downloadCert(order)
}
func ShowMenu(){
	fmt.Println("========欢迎使用环智中诚™ · 非官方,MJJ客户端===============")
	fmt.Println("1. 签新证书(域名/IP)")
	fmt.Println("2. 下载证书")
	fmt.Println("3. 生成CSR,无需提前生成，签证书会自动调用。")
	fmt.Println("4. 查看订单信息")
	fmt.Println("5. 提交订单（由于没有获取所有订单接口，只能手工提交）")
	fmt.Println("6. 查看可签的证书列表")
	fmt.Println("7. 自动续签配置")
	fmt.Println("================非以上选项[1-7]将退出程序=================")
	sel := ComInputNum("请选择")
	if sel == 1 {
		newCert()
	}
	if sel == 2 {
		downloadCert(nil)
	}
	if sel == 3 {
		genCSR(true)
	}
	if sel == 4 {
		showOrders()
	}
	if sel == 5 {
		submitOrder()
	}
	if sel == 6 {
		showCertList()
	}
	if sel == 7 {
		fmt.Printf("\r\n\r\n\r\n\r\n---->该功能还没开发，敬请期待<----\r\n\r\n\r\n\r\n")
		ShowMenu()
	}
	os.Exit(0)
}

func ShowEntry(){
	fmt.Println("=======环智中诚™ (TrustOcean Limited)====")
	fmt.Println("-->本程序非官方客户端，使用本程序所造成的任何法律责任，由使用者承担")
	fmt.Println("-->根据https://github.com/londry/Encryption365_Baota 以及相关文档而开发")
	fmt.Println("-->官方授权协议：")
	fmt.Println(`Copyright (c) 2019 TrustOcean Limited
   Permission is hereby granted, free of charge,
   to any person obtaining a copy of this software and associated documentation files (the "Software"),
   to deal in the Software without restriction,
   including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
   and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
   subject to the following conditions:
   The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
   INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
   IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
   WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`)
	fmt.Println("========================================")
	if !ComSelect("是否同意") {
		os.Exit(0)
	}
}