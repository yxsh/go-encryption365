package req

import (
	"encoding/json"
	"errors"
	"fmt"
)
var ProductIndex =make(map[string]*Product365)
type Product365 struct {
	Id string
	Title string
	//使用范围
	UseAge string
	IsFree bool
	Level string
	Period string
	PeriodText string
	//单一价格
	Fqdn float64
	//宽字价格
	WildCard float64
}
func ProductFromJSon(jsonStr string) ([]*Product365, error){
	if jsonStr == "" {
		return nil,errors.New("jsonStr empty")
	}
	data := make(map[string]interface{})
	e := json.Unmarshal([]byte(jsonStr), &data)
	if e != nil {
		return nil, e
	}
    ps := data["products"].(map[string]interface{})
    if ps == nil  || len(ps) < 1{
    	return nil, errors.New("not any more product")
	}
	products := make([]*Product365, len(ps))
	i := 0
	for key,val := range ps {
		mapVal := val.(map[string]interface{})
		product := new(Product365)
		product.Id = key
		product.Title = mapVal["title"].(string)
        product.Period = mapVal["period"].(string)
        product.PeriodText = mapVal["periodText"].(string)
        product.IsFree = mapVal["isFree"].(bool)
        product.Level = mapVal["level"].(string)
        product.UseAge = mapVal["useage"].(string)
        price := mapVal["price"].(map[string]interface{})
        product.Fqdn = price["fqdn"].(float64)
        product.WildCard = price["wildcard"].(float64)
		products[i] = product
		ProductIndex[key] = product
		i++
	}

	return products, nil
}
func (p *Product365) String() string{
	fmtStr := `
==============================================
      产品编号: %s
      名称:    %s 
      授权方式: %s
      是否免费: %s
      证书级别：%s
      付款周期：%s
      单一价格：%.2f
      泛域价格：%.2f
==============================================
`
freeText := "是"
if !p.IsFree {
	freeText = "否"
}
return fmt.Sprintf(fmtStr,
	p.Id,
	p.Title,
	p.UseAge,
	freeText,
	p.Level,
	p.PeriodText,
	p.Fqdn,
	p.WildCard,
	)
}
