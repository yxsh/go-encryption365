package req

import (
   "encoding/json"
   "errors"
)

type ReqCertInfo struct {
   ProductId string
   Period string
   CsrCode string
   Domains string `desc:"IP/域名， 多个请用英文(,)隔开(泛域名填 *.xxx.com)"`
}
type CertInfo365 struct {
   CertStatus string `json:"cert_status"`
   CreatedAt string  `json:"created_at"`
   TrustoceanId int64 `json:"trustocean_id"`
   CsrCode string `json:"csr_code"`
}
func CertFromJson(jsonStr string) (*CertInfo365, error){
  if  jsonStr == "" {
     return nil, errors.New("jsonStr Empty")
  }
  info := new(CertInfo365)
  err := json.Unmarshal([]byte(jsonStr), info)
  if err != nil {
     return nil, err
  }
  return info, nil
}