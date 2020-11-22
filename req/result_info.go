package req

import "fmt"

type ResultInfo365 struct {
	Status string `json:"result"`
	Message string `json:"message"`
	OrgText string
}
func (r *ResultInfo365) IsOK() bool{
	return "success" == r.Status
}
func (r *ResultInfo365) String() string{
	return fmt.Sprintf("%s-%s:%s", r.Status, r.Message, r.OrgText)
}