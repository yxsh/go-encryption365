package ui

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func ComSelect(text string) bool{
	tmp := ComInput(text + "[y/n]")
	return tmp == "y" || tmp == "Y"
}
func ComInput(text string) string{
	var tmp string
	fmt.Print(text, ":")
	_,e := fmt.Scanln(&tmp)
	if e != nil {
		fmt.Scan(&tmp)
	}
	return tmp
}
func ComInputNum(text string) int64{
	str := ComInput(text)
	if str == "" {
		return -1
	}
	i, e := strconv.ParseInt(str,10, 64)
	if e != nil {
		return -1
	}
	return i
}
func ComLongText(container string) string{
   tmp := ""
   read := bufio.NewReader(os.Stdin)
   for !strings.Contains(tmp, container) {
   	 b, e := read.ReadByte()
   	 if e != nil{
   	 	return tmp
	 }
	 tmp += string(b)
   }
   return tmp
}
func ComAnyKeyContinue(){
	fmt.Print("\r\n按任意键继续...")
	read := bufio.NewReader(os.Stdin)
	read.ReadByte()
}
func InputPkg(st interface{})  interface{}{
	if st == nil{
		return nil
	}

	val := reflect.ValueOf(st)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	n := val.NumField()
	for i := 0 ; i < n ; i++ {
		filed := val.Type().Field(i)
		tag := filed.Tag
		desc := tag.Get("desc")
		if  desc == "" {
			continue
		}
		showDesc := "请输入" + desc
		if filed.Type.Kind() == reflect.Int {
			val.FieldByName(filed.Name).SetInt(ComInputNum(showDesc))
			continue
		}
		f := val.FieldByName(filed.Name)
		str := f.String()
		if str == "" {
			f.SetString(ComInput(showDesc))
		}
	}
	return st
}