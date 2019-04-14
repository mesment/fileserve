package setting

import (
	"github.com/gpmgo/gopm/modules/log"
	"github.com/mesment/personblog/pkg/initconfig"
	"io/ioutil"
)

//配置信息结构体
type Config struct {
	AppSetting APP  	`ini:"app"`
	ServSetting Server	`ini:"server"`
	DbSetting DataBase 	`ini:"database"`
	Log 	Log 		`ini:"log"`
}

//对应ini配置文件中app配置
type APP struct {
	JwtSecret string	`init:"jwt_secret"`	//JWT 加密字符串
	PageSize  int		`ini:"page_size"`	//每页显示文章条数
}

//对应ini log配置
type Log struct {
	LogLevel	string		`ini:"level"`
	LogPath  	string 		`ini:"log_path"` 	//日志目录
	LogPrefix 	string		`ini:"log_prefix"`	//日志前缀
	LogExt 		string		`ini:"log_ext"`		//日志后缀
	TimeFormat	string		`ini:"time_format"`	//日志文件名日期格式
	SplitSize int64			`ini:"split_size"`
}

//对应ini配置文件中server配置
type Server struct {
	RunModel  	string	`ini:"run_mode""`
	HttpPort	int		`ini:"http_port"`
	ReadTimeOut  int	`ini:"read_timeout"`
	WriteTimeOut int	`ini:"write_timeout"`
}

//对应ini配置文件中database配置
type DataBase struct {
	Type 	string  	`ini:"type"`		//数据库类型mysql
	User 	string		`ini:"user"` 		//用户名
	PassWd	string		`ini:"password"` 	//密码
	Name 	string		`ini:"name"`		//数据库名
	Host 	string		`ini:"host"` 		//主机地址 127.0.0.1:3306
}

//配置信息
var Cfg Config

func SetUp()  {

	data, err := ioutil.ReadFile("config/app.ini")
	if err != nil {
		log.Error("读取配置文件失败：%v",err)
	}

	err = iniconfig.UnMarshal(data, &Cfg)
	if err != nil {
		log.Error("unmarshal failed, err:%v", err)
		return
	}

	/*
	var err error
	cfg, err = ini.Load("config/app.ini")
	fmt.Printf("cfg: %v\n",cfg)
	if err != nil {
		log.Fatalf("解析配置文件'config/app.ini'失败：%s",err.Error())

	}
	//映射app 配置
	mapTo("app", AppSetting)
	log.Printf("AppSetting：%v",AppSetting)
	//映射server配置
	mapTo("server", ServSetting)
	log.Printf("serverSetting：%v",ServSetting)
	//映射database配置
	mapTo("database",DbSetting)
	log.Printf("执行完dbsetting")
	log.Printf("dbSetting：%v",DbSetting)
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s Setting err: %v", section,err)
	}


	*/
}

