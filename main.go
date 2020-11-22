package main

import (
	"github.com/yxsh/go-encryption365/conf"
	"github.com/yxsh/go-encryption365/presenter"
	"github.com/yxsh/go-encryption365/ui"
)

func main(){
	ui.Presenter = presenter.New()
	ui.Presenter.Config = conf.Load()
	ui.ShowEntry()
	if ui.Presenter.Config.UserName != "" {
		ui.ShowTip()
	}else{
		ui.ShowLoginTip()
	}
}
