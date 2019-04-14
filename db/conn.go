package db

import (
	"database/sql"
	"fmt"
	"github.com/mesment/fileserver/pkg/setting"
	"log"
	"net/url"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB



func Setup()  {

	var dbtype = setting.DBCfg.DBType
	var user = setting.DBCfg.User
	var passwd = setting.DBCfg.Password
	var host = setting.DBCfg.Host
	var dbname = setting.DBCfg.DBName

	timezone := "'Asia/Shanghai'"  //设置时区，mysql默认是utc时间
	fmtstr := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&time_zone="

	//db, err := sql.Open("mysql", "user:password@/dbname?charset=utf8mb4&parseTime=true&loc=Local&time_zone=" + url.QueryEscape(timezone))
	dns := fmt.Sprintf(fmtstr,user,passwd,host,dbname) + url.QueryEscape(timezone)

	uri :=fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=%s",
		user,
		passwd,
		host,
		dbname,
		url.QueryEscape("Asia/Shanghai" ),
		)

	log.Printf("dns:%s\n",dns)
	log.Printf("uri:%s\n",uri)

	db, _ = sql.Open(dbtype,uri)
	db.SetMaxOpenConns(1000)

	err := db.Ping()
	if err != nil {
		log.Printf("连接数据库失败:%s", err.Error())
		os.Exit(1)
	}
}

//返回数据库连接对象
func DBConn() *sql.DB  {
	return db;
}
