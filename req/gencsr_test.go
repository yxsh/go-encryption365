package req

import (
	"fmt"
	"testing"
)

func TestCSR365_Generate(t *testing.T){
	csr := NewCsr()
	err := csr.Generate()
	if err != nil{
		t.FailNow()
	}
	fmt.Println(csr.PublicCer)
	fmt.Println(csr.PrivateKey)
}