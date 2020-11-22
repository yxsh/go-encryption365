package req

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)
type DcvInfo365 struct {
	Domain string
	Emails []string
	DnsHost string
	DnsType string
	DnsValue string
	HttpVerifyLink string
	HttpsVerifyLink string
	FileName string
	FileContent string
}
type OrderInfo365 struct {
	Result string `json:"result"`
	TrustoceanId int `json:"trustocean_id"`
	CertStatus string `json:"cert_status"`
	CertStatusText string
	ProductName string `json:"product_name"`
	DomainCount int `json:"domain_count"`
	CertCode string `json:"cert_code"`
	CreatedAt string `json:"created_at"`
	ExpireAt string `json:"expire_at"`
	Domains []string
	DcvInfo map[string]interface{} `json:"dcv_info"`
	DcvInfos []*DcvInfo365
	DcvIndex map[string]*DcvInfo365
}
type OrderStatus map[string]string

func (o *OrderInfo365) IsOk() bool{
	return "success" == o.Result
}
func (o *OrderInfo365) SaveText() error{
	for _, d := range o.DcvIndex {
		file, e := os.Create(d.FileName)
		if e != nil {
			return e
		}
		_, e = file.WriteString(d.FileContent)
		if e != nil {
			return e
		}
		file.Close()
	}
	return nil
}
func OrderFromJSON(jsonStr string) (*OrderInfo365, error){
	info := new(OrderInfo365)
	info.DcvIndex = make(map[string]*DcvInfo365)
	err := json.Unmarshal([]byte(jsonStr), info)
	if err != nil {
		return nil, err
	}
	if info.DcvInfo != nil {
		domains := make([]string, len(info.DcvInfo))
		infos := make([]*DcvInfo365, len(info.DcvInfo))
		i := 0
		for key,val := range info.DcvInfo {
			domains[i] = key
			dcv := val.(map[string]interface{})
			es :=  dcv["emails"].([]interface{})
			emails := make([]string, 1)
			if es != nil {
				for _,d := range es {
					emails = append(emails, d.(string))
				}
			}
			info365 := &DcvInfo365{
				Domain: dcv["domain"].(string),
				Emails: emails,
				DnsHost: dcv["dns_host"].(string),
				DnsType: dcv["dns_type"].(string),
				DnsValue: dcv["dns_value"].(string),
				HttpVerifyLink: dcv["http_verifylink"].(string),
				HttpsVerifyLink: dcv["https_verifylink"].(string),
				FileName: dcv["https_filename"].(string),
				FileContent: dcv["https_filecontent"].(string),
			}
			infos[i] = info365
			info.DcvIndex[key] = info365
			i++
		}
	}
	info.CertStatusText = NewOrderStatus().Get(info.CertStatus)
	return info, nil
}
func (o *OrderInfo365) ShowDomain() string{
	if o.CertStatus != "enroll_caprocessing" || len(o.DcvIndex) < 1{
		return ""
	}
	ret := "-----------------待验证-------------------------"
	for _, d := range o.DcvIndex {
		fmtStr := `
       验证域名所有权，可以选择 文件验证，DNS验证中的一项.(IP类只能文件验证)
       ---> 相关验证文件已生成到当前目录

       域名： %s
       DNS验证类型: %s
       DNS主机记录：%s
       DNS记录值:   %s

       文件名: %s
       http路径如下: %s
       https路径如下：%s
       文件内容:
        %s
    `
		ret += fmt.Sprintf(fmtStr,
			d.Domain,
			d.DnsType,
			d.DnsHost,
			d.DnsValue,
			d.FileName,
			d.HttpVerifyLink,
			d.HttpsVerifyLink,
			d.FileContent,
		)
	}
	ret += "---------------------------------------------"
	return ret
}
func (o *OrderInfo365) String() string{
	fmtStr := `
==================================================
    订单编号: %d
    状态： %s
    产品名称：%s
    域名数：%d
    IP/域名: %s
    创建时间: %s
    到期时间: %s
    %s
    %s
==================================================    
`
e := o.SaveText()
errStr := ""
if e != nil {
	errStr = e.Error()
}
return fmt.Sprintf(fmtStr, o.TrustoceanId, o.CertStatusText,
	o.ProductName, o.DomainCount, strings.Join(o.Domains, ","),
	o.CreatedAt, o.ExpireAt, o.ShowDomain(), errStr)
}
func NewOrderStatus() OrderStatus  {
	status := make(OrderStatus, 6)
	status["issued_active"] = "签发完成"
	status["enroll_caprocessing"] = "签发中(待验证)"
	status["expired"] = "已过期"
	status["cancelled"] = "已取消"
	status["revoked"] = "已被吊销"
	status["rejected"] = "SSL订单已被CA机构拒绝"
	return status
}
func (o OrderStatus) Get(status string) string{
	stat := o[status]
	if stat == "" {
		return "未知状态"
	}
	return stat
}
