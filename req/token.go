package req

import (
	"encoding/json"
	"fmt"
)

type Token365 struct {
	ClientId    string `json:"client_id"`
	AccessToken string `json:"access_token"`
}
func NewTokenFromJson(jsonStr string)*Token365{
	if jsonStr == "" {
		return nil
	}
	token365 := new(Token365)
	err := json.Unmarshal([]byte(jsonStr), token365)
	if err != nil {
		return nil
	}
	return token365
}
func (r *Token365)String() string{
	return fmt.Sprintf("client_id=%s  ,  access_token=%s", r.ClientId, r.AccessToken)
}