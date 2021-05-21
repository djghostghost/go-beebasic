package b_logger

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type FileConfig struct {
	FileName string   `json:"filename"`
	MaxLines int      `json:"maxlines"` // 每个文件保存的最大行数，默认值 1000000
	MaxSize  int      `json:"maxsize"`  // 每个文件保存的最大尺寸, 默认值是 1 << 28
	Daily    bool     `json:"daily"`    // 是否按照每天 logrotate，默认是 true
	MaxDays  int      `json:"maxdays"`  // 文件最多保存多少天，默认保存 7 天
	Hourly   bool     `json:"hourly"`
	MaxHours int64    `json:"maxhours"`
	Rotate   bool     `json:"rotate"`   // 默认是 true
	Level    int      `json:"level"`    // 默认是 Trace 级别
	Perm     string   `json:"perm"`     // 日志文件权限
	Separate []string `json:"separate"` // 需要单独写入文件的日志级别,设置后命名类似 test.error.log
}

type AllConfig struct {
	FileLoggers map[string]FileConfig `toml:"file"`
}

var loggerConfig AllConfig
var lock sync.Mutex

var AppConfig FileConfig

func Init() {
	readConfig()
	AppConfig = loggerConfig.FileLoggers["app"]
	content, _ := json.Marshal(AppConfig)
	log.Printf("logs config %s\n", string(content))
	if AppConfig.Separate == nil {
		err := logs.SetLogger("file", string(content))
		if err != nil {
			panic(err)
		}

	} else {
		err := logs.SetLogger("multifile", string(content))
		if err != nil {
			panic(err)
		}
	}
	logs.SetLogFuncCall(true)
	// 输出access日志
	web.BConfig.Log.AccessLogs = true
	// orm 日志也输入到系统日志中
	orm.DebugLog = orm.NewLog(logs.GetBeeLogger())
}

func BuildCustomLogger(loggerName string) *logs.BeeLogger {
	loggerFileConfig := loggerConfig.FileLoggers[loggerName]
	configContent, _ := json.Marshal(loggerFileConfig)
	result := logs.NewLogger()
	if loggerFileConfig.Separate == nil {
		result.SetLogger("file", string(configContent))
	} else {
		result.SetLogger("multifile", string(configContent))
	}
	return result
}

// read config
func readConfig() {
	lock.Lock()
	defer lock.Unlock()
	filename := "logger.toml"
	runmode := web.AppConfig.DefaultString("runmode","dev")
	_, err := os.Stat("./conf/" + runmode + "." + filename)
	if err == nil {
		filename = runmode + "." + filename
	}

	if len(loggerConfig.FileLoggers) == 0 {
		data, err := ioutil.ReadFile("./conf/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		var loggerToml AllConfig
		if _, err := toml.Decode(string(data), &loggerToml); err != nil {
			log.Fatal(err)
		}
		loggerConfig = loggerToml
		for _, fileLogger := range loggerConfig.FileLoggers {
			if fileLogger.MaxDays == 0 {
				fileLogger.MaxDays = 7
			}
			if fileLogger.MaxLines == 0 {
				fileLogger.MaxLines = 100000
			}
			if fileLogger.MaxSize == 0 {
				fileLogger.MaxSize = 1 << 28
			}
		}
	}
}
