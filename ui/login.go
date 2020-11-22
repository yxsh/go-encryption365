package ui

import (
	"fmt"
	"github.com/yxsh/go-encryption365/req"
	"os"
)
func register(userName string){
	info := new(req.RegisterInfo)
	info.UserName = userName
    InputPkg(info)
	if !Presenter.Register(info) {
		if !ComSelect("是否重新填写注册？") {
			os.Exit(0)
		}
		register(userName)
	}
    login(info.UserName, info.Password)
}
func login(user, pwd string){
   if user == "" || pwd == "" {
   	 user = ComInput("请输入用户名")
   	 pwd = ComInput("请输入密码")
   }
   fmt.Printf("现在以\r\n用户名：%s \r\n 密码: %s \r\n登陆--->\r\n", user, pwd)
   if !Presenter.Login(user, pwd) {
   	  if !ComSelect("是否重新登陆?") {
   	     os.Exit(0)
	  }
	  login("", "")
   }

}
func sendCode(){
	email := ComInput("请输入注册邮件(用于接收注册码)")
	if !Presenter.SendCode(email) {
		if !ComSelect("发送失败,是否重发") {
			os.Exit(0)
		}
		sendCode()
	}
	fmt.Println("是否已收到邮件? 下一步将需要验证码")
	fmt.Println("==============================")
	fmt.Println("1.没有收到,重新发送")
	fmt.Println("2.已收到，继续")
	fmt.Println("==============================")
	sel := ComInputNum("请选择")
	if sel == 1 {
		sendCode()
	}
	if sel == 2 {
		register(email)
	}
}
func ShowTip(){
	fmt.Println("=====系统已检测本地配置文件=========")
	fmt.Printf("1.加载配置文件,并继续(%s)\r\n", Presenter.Config.UserName)
	fmt.Println("2.放弃配置文件，初始化")
	fmt.Println("==============================")
	sel := ComInputNum("请选择")
	if sel == 1 {
		Presenter.RefreshToken()
		ShowMenu()
	}
	if sel == 2 {
		ShowLoginTip()
		os.Exit(0)
	}
}
func ShowLoginTip(){
	fmt.Println("===>控制台地址:https://console.trustocean.com/certificate/create/91")
	fmt.Println("===》环智中诚™账号")
	fmt.Println(" 1. 新用户注册")
	fmt.Println(" 2. 已有账号登陆")
	sel := ComInputNum("请选择")
	if sel != 1 && sel != 2 {
		os.Exit(0)
	}
    if sel == 1 {
		sendCode()
	}
    if sel == 2 {
    	login("", "")
	}
}

