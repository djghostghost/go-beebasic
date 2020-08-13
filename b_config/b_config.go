package b_config

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/zouyx/agollo"
)

type configEvent struct {
	Namespace string
	Changes   map[string]changeDate
}
type changeDate struct {
	OldValue   string
	NewValue   string
	ChangeType int
}

var sysConfigKey = map[string]struct{}{
	"AppName":                {},
	"RunMode":                {},
	"RouterCaseSensitive":    {},
	"ServerName":             {},
	"RecoverPanic":           {},
	"RecoverFunc":            {},
	"CopyRequestBody":        {},
	"EnableGzip":             {},
	"MaxMemory":              {},
	"EnableErrorsShow":       {},
	"EnableErrorsRender":     {},
	"Graceful":               {},
	"ServerTimeOut":          {},
	"ListenTCP4":             {},
	"EnableHTTP":             {},
	"AutoTLS":                {},
	"Domains":                {},
	"TLSCacheDir":            {},
	"HTTPAddr":               {},
	"HTTPPort":               {},
	"EnableHTTPS":            {},
	"HTTPSAddr":              {},
	"HTTPSPort":              {},
	"HTTPSCertFile":          {},
	"HTTPSKeyFile":           {},
	"EnableAdmin":            {},
	"AdminAddr":              {},
	"AdminPort":              {},
	"EnableFcgi":             {},
	"EnableStdIo":            {},
	"AutoRender":             {},
	"EnableDocs":             {},
	"FlashName":              {},
	"FlashSeparator":         {},
	"DirectoryIndex":         {},
	"StaticDir":              {},
	"StaticExtensionsToGzip": {},
	"TemplateLeft":           {},
	"TemplateRight":          {},
	"ViewsPath":              {},
	"EnableXSRF":             {},
	"XSRFKey":                {},
	"XSRFExpire":             {},
	"AccessLogs":             {},
	"EnableStaticLogs":       {},
	"AccessLogsFormat":       {},
	"FileLineNum":            {},
}

func init() {
	readyConfig := &agollo.AppConfig{
		AppId:         beego.AppConfig.String("apollo.appid"),
		Cluster:       beego.AppConfig.String("apollo.cluster"),
		NamespaceName: beego.AppConfig.String("apollo.namespace"),
		Ip:            beego.AppConfig.String("apollo.mateserver"),
	}

	agollo.InitCustomConfig(func() (*agollo.AppConfig, error) {
		return readyConfig, nil
	})
}

//设置业务参数
func setBizConfig() {
	event := agollo.ListenChangeEvent()
	for changeEvent := range event {
		bytes, _ := json.Marshal(changeEvent)
		var result configEvent
		err := json.Unmarshal(bytes, &result)
		if err != nil {
			beego.Error("[Apollo] Unmarshal Error :%s]", err.Error())
		}
		for k, v := range result.Changes {
			switch v.ChangeType {
			case 0:
				beego.Info(fmt.Sprintf("[Apollo] BizConfig Add Key:%s, OldValue:%s,NewValue:%s",
					k, v.OldValue, v.NewValue))
			case 1:
				beego.Info(fmt.Sprintf("[Apollo] BizConfig Modify Key:%s, OldValue:%s,NewValue:%s",
					k, v.OldValue, v.NewValue))
			case 2:
				beego.Info(fmt.Sprintf("[Apollo] BizConfig Del Key:%s, OldValue:%s,NewValue:%s",
					k, v.OldValue, v.NewValue))
			}
			_, isSysConfig := sysConfigKey[k]
			if !isSysConfig {
				err := beego.AppConfig.Set(k, v.NewValue)
				if err != nil {
					beego.Error("[Apollo] Set Config Error :%s,Key:%s, OldValue:%s,NewValue:%s]", err.Error(),
						k, v.OldValue, v.NewValue)
				}
			}
		}
	}
}

//设置容器参数
func setApplicationConfig() {
	// App Config
	RouterCaseSensitive := agollo.GetBoolValue("RouterCaseSensitive", false)
	beego.BConfig.RouterCaseSensitive = RouterCaseSensitive

	ServerName := agollo.GetStringValue("ServerName", "beego")
	beego.BConfig.ServerName = ServerName

	RecoverPanic := agollo.GetBoolValue("RecoverPanic", true)
	beego.BConfig.RecoverPanic = RecoverPanic

	CopyRequestBody := agollo.GetBoolValue("CopyRequestBody", true)
	beego.BConfig.CopyRequestBody = CopyRequestBody

	EnableErrorsShow := agollo.GetBoolValue("EnableErrorsShow", true)
	beego.BConfig.EnableErrorsShow = EnableErrorsShow

	EnableErrorsRender := agollo.GetBoolValue("EnableErrorsRender", true)
	beego.BConfig.EnableErrorsRender = EnableErrorsRender

	//Web Config
	AutoRender := agollo.GetBoolValue("AutoRender", true)
	beego.BConfig.WebConfig.AutoRender = AutoRender

	EnableDocs := agollo.GetBoolValue("EnableDocs", false)
	beego.BConfig.WebConfig.EnableDocs = EnableDocs

	FlashName := agollo.GetStringValue("FlashName", "BEEGO_FLASH")
	beego.BConfig.WebConfig.FlashName = FlashName

	FlashSeparator := agollo.GetStringValue("FlashSeparator", "BEEGOFLASH")
	beego.BConfig.WebConfig.FlashSeparator = FlashSeparator

	DirectoryIndex := agollo.GetBoolValue("DirectoryIndex", false)
	beego.BConfig.WebConfig.DirectoryIndex = DirectoryIndex
	//Listen Config

	HTTPPort := agollo.GetIntValue("HTTPPort", 8080)
	beego.BConfig.Listen.HTTPPort = HTTPPort

	HTTPAddr := agollo.GetStringValue("HTTPAddr", "localhost")
	beego.BConfig.Listen.HTTPAddr = HTTPAddr

	EnableAdmin := agollo.GetBoolValue("EnableAdmin", false)
	beego.BConfig.Listen.EnableAdmin = EnableAdmin

	AdminAddr := agollo.GetStringValue("AdminAddr", "localhost")
	beego.BConfig.Listen.AdminAddr = AdminAddr

	AdminPort := agollo.GetIntValue("AdminPort", 8090)
	beego.BConfig.Listen.HTTPPort = AdminPort
	//Log Config
	AccessLogs := agollo.GetBoolValue("AccessLogs", false)
	beego.BConfig.Log.AccessLogs = AccessLogs
	FileLineNum := agollo.GetBoolValue("FileLineNum", true)
	beego.BConfig.Log.FileLineNum = FileLineNum
}

func RegisterConfig() error {
	err := agollo.Start()
	if err != nil {
		beego.Error("[Apollo] Start Error. :%s]", err.Error())
		return err
	}
	setApplicationConfig()
	go setBizConfig()
	return nil
}
