package setting

import (
	"gopkg.in/ini.v1"
	"log"
	"strings"
)

var (
	Cfg *ini.File
	RunModel string
	DBType string
	User string
	Password string
	Host  string
	DBName string

)


type DataBase struct {
	DBType string
	User string
	Password string
	Host  string
	DBName string
}

var DBCfg  = &DataBase{}

type Server struct {
	JwtSecret string
	FileStorePath string

}

var ServerCfg  = &Server{}

func Setup()  {
	var err error
	Cfg,err =ini.LooseLoad("config/fileserver.ini","config/server.ini")
	if err != nil {
		log.Fatalf("Failed to parse config server.ini.%s",err.Error())
		return
	}

	err = Cfg.Section("database").MapTo(DBCfg)
	if err != nil {
		log.Printf("setting 映射配置文件database失败：%s",err)
		return
	}

	err = Cfg.Section("server").MapTo(ServerCfg)
	if err != nil {
		log.Printf("setting 映射配置文件server失败：%s",err)
		return
	}

	//判断文件服务的存储路径结尾是否包含/，如果有去掉路径最后的/
	fileStorePath := ServerCfg.FileStorePath
	filePathLen := len(fileStorePath)
	idx := strings.LastIndex(fileStorePath,"/")
	if len(fileStorePath) == 0 {
		log.Fatal("配置信息中文件存储路径不能为空")
	}
	if  idx == filePathLen -1 {
		//去掉结尾的/
		fileStorePath = fileStorePath[:filePathLen - 1 ]
	}
	ServerCfg.FileStorePath = fileStorePath

	//LoadDataBase()
}

func LoadDataBase()  {
	section, err := Cfg.GetSection("database")
	if err != nil {
		log.Fatalln("解析配置文件database 失败：%s",err.Error())
		return
	}
	DBType = section.Key("DBType").MustString("mysql")
	User = section.Key("User").String()
	Password = section.Key("Password").String()
	Host = section.Key("Host").String()
	DBName = section.Key("DBName").String()

}
