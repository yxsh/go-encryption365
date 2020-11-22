package req

type RegisterInfo struct {
	UserName string `json:"username" desc:"用户名(邮件)"`
	Password string `json:"password" desc:"密码"`
	AuthCode string `json:"authcode" desc:"邮件验证码"`
	RealName string `json:"realName" desc:"真实姓名"`
	IdCardNumber string `json:"idcardNumber" desc:"身份证号"`
	PhoneNumber string `json:"phoneNumber" desc:"手机号"`
	Country string `json:"country" desc:"国家代码(CN - 为中国)"`
	CompanyName string `json:"companyname" desc:"公司名称"`
}
func (r *RegisterInfo) toValues() Value{
   return StructToValue(r)
}
