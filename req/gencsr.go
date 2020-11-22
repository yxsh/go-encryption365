package req

import (
  "errors"
  "io/ioutil"
  "net/http"
  "net/url"
  "strings"
)

/**
C=CN&ST=ShangHai&L=Shanghai&O=Youtudata&OU=com&CN=127.0.0.1&keySize=2048
application/x-www-form-urlencoded
 */
type CSR365 struct {
  Country string `desc:"国家编号（中国请输入CN）"`
  StateOrProvinceName string `desc:"省或者直辖市"`
  LocalityName string `desc:"地区"`
  OrganizationName string `desc:"公司"`
  OrganizationalUnitName string `desc:"部门"`
  //IP证书默认写入: common-name-for-public-ip-address.com
  CommonName string `desc:"域名"`
  KeySize string `desc:"加密强度：2048，4096"`
  PublicCer string
  PrivateKey string
}

const GenUrl  = "https://www.csr.sh/generate"
func NewCsr()*CSR365{
   return &CSR365{
     Country:                "CN",
     StateOrProvinceName:    "Xian",
     LocalityName:           "Shaanxi",
     OrganizationName:       "Encryption365 SSL Security",
     OrganizationalUnitName: "Encryption365 SSL Security",
     CommonName:             "common-name-for-public-ip-address.com",
     KeySize:                "2048",
     PublicCer:              "",
     PrivateKey:             "",
   } 
}
func (c *CSR365) parse(s string) error{
  index := strings.Index(s, "-----BEGIN PRIVATE KEY-----")
  if index < 0 {
    return errors.New("ErrorResponse")
  }
  c.PublicCer = s[:index]
  c.PrivateKey = s[index:]
  return nil
}
func (c *CSR365)Generate() error{
  values := make(url.Values)
  values.Add("C", c.Country)
  values.Add("ST", c.StateOrProvinceName)
  values.Add("L", c.LocalityName)
  values.Add("O", c.OrganizationName)
  values.Add("OU", c.OrganizationalUnitName)
  values.Add("CN", c.CommonName)
  values.Add("keySize", c.KeySize)
  rsp,err := http.PostForm(GenUrl, values)
  if err != nil {
    return err
  }
  if rsp.StatusCode != 200 {
    return errors.New("ServerError")
  }
  result,err := ioutil.ReadAll(rsp.Body)
  if err != nil {
    return err
  }
  return c.parse(string(result))
}