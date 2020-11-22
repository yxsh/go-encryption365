package req
/**
{
  "title": "Encryption365™ 免费证书",
  "tip":"lib",
  "type":"扩展",
  "php": "index.php",
  "name": "encryption365",
  "ps": "更专业的 SSL 安全证书管理客户端, 免费自动化部署 SSL 证书, 可免费保护高达1000条域名和公网IPv4地址。支持大规模部署和集中化管理。",
  "versions": "1.2",
  "checks": "/www/server/panel/plugin/encryption365",
  "author": "环智中诚™",
  "date":"2020-05-01",
  "home": "https://www.trustocean.com",
  "default":false,
  "display":1
}
文档:
https://support.trustocean.com/doc/C4WDG8Zri0/API+%E6%96%87%E6%A1%A3-f9K6pRVKsm
 */
import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type Request365 struct {
	BaseURL string
	Client http.Client
	Token *Token365
}
type Value map[string]string
const (
	 UserAgent = "Encryption365-Client/1.2;BaotaPanel-LinuxVersion"
	 baseURL = "https://encrypt365.trustocean.com"
	 contextType = "application/json"
	)
func NewClient() *Request365 {
	jar,_ := cookiejar.New(nil)
	ts := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Request365{
		BaseURL: baseURL,
		Client: http.Client{
			Jar:jar,
			Transport: ts,
			Timeout: time.Second * 40,
		},
	}
}
func (r *Request365)SetToken(clientId, token string){
   r.Token = &Token365{
	   ClientId:   clientId,
	   AccessToken: token,
   }
}
func (r *Request365) doRequest(url string, values Value)(info365 *ResultInfo365, err error){
	if values == nil {
		values = make(Value)
	}
	if r.Token != nil {
		values["client_id"] = r.Token.ClientId
		values["access_token"] = r.Token.AccessToken
	}
	v, _ := json.Marshal(values)
	reqBody := strings.NewReader(string(v))
	req365,err := http.NewRequest("POST", r.BaseURL + url, reqBody)
	if err != nil {
		return nil, err
	}
	req365.Header.Add("Content-Type", contextType)
	req365.Header.Add("User-Agent", UserAgent)
	rsp,err := r.Client.Do(req365)
	if err != nil {
		return nil, err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil,err
	}
	if rspBody == nil  || len(rspBody) < 1{
		return nil, errors.New("empty result")
	}
	info := new(ResultInfo365)
	err = json.Unmarshal(rspBody, info)
	if err != nil {
		return nil, err
	}
	info.OrgText = string(rspBody)
    return info, nil
}
func (r *Request365) SendAuthCode(userName string)(info365 *ResultInfo365, err error){
    param := make(Value)
    param["username"]= userName
	ret, e := r.doRequest("/account/register/authcode", param)
	if e != nil {
		return nil,e
	}
	return  ret,nil
}

func (r *Request365)AccountRegister(info *RegisterInfo)(info365 *ResultInfo365, err error){
	ret,e := r.doRequest("/account/register", info.toValues())
	if e != nil {
		return nil,e
	}
	return  ret,nil
}
func (r *Request365)ClientCreate(userName,Password string) error{
	formData := make(Value)
	formData["username"] = userName
	formData["password"] = Password
	formData["servername"] = "example.com"
	ret, e := r.doRequest("/client/create", formData)
    if e != nil{
    	return e
	}
	if !ret.IsOK() {
		return errors.New(ret.Message)
	}
	userToken := NewTokenFromJson(ret.OrgText)
	if userToken == nil {
		return errors.New("get token fail")
	}
	r.Token = userToken
	return nil
}
func (r *Request365) CreateNewCert(info *ReqCertInfo)(*OrderInfo365, error){
	param := make(Value)
	param["pid"] = info.ProductId
	param["period"] = info.Period
	param["csr_code"] = info.CsrCode
	param["domains"] = info.Domains
	param["renew"] = "false"
	param["old_vendor_id"] = "-1"
	ret, e := r.doRequest("/cert/create", param)
	if e != nil{
		return nil, e
	}
	if !ret.IsOK() {
		return nil, errors.New(ret.Message)
	}
	c, e := OrderFromJSON(ret.OrgText)
	if e != nil {
		return nil, e
	}
	return c, nil
}
func (r *Request365) GetOrderInfo(id string) (*OrderInfo365, error){
	param := make(Value)
	param["trustocean_id"] = id
	ret, e := r.doRequest("/cert/details", param)
	if e != nil {
		return nil, e
	}
	if !ret.IsOK() {
		return nil, errors.New(ret.Message)
	}
	order, e := OrderFromJSON(ret.OrgText)
	if e != nil {
		return nil, e
	}
	return order, nil
}
func (r *Request365) GetProducts()([]*Product365, error){
	ret, e := r.doRequest("/account/products", nil)
	if e != nil{
		return nil, e
	}
    if !ret.IsOK() {
		return nil, errors.New(ret.Message)
	}
	p, e := ProductFromJSon(ret.OrgText)
	if e != nil {
		return nil, e
	}
	return p, nil
}