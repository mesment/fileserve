package db

import (
	"log"
	"time"
)

type UserFileTable struct {
	UserName string
	FileHash	string
	FileMD5		string
	FileName 	string
	FileSize    int64
	FileAddr    string
	UploadAt	string
	LastUpdated string
}

//保存文件到用户信息表
func SaveUserFileInfo(username,filename,sha1, md5, fileAddr string, fileSize int64) (bool,error) {
	stmt, err := DBConn().Prepare(`insert  into user_file(user_name,file_sha1,file_md5,file_name,
										file_size,file_addr) values (?,?,?,?,?,?)`)
	if err != nil {
		return false,err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username,sha1,md5,filename,fileSize,fileAddr)
	if err != nil {
		return false, err
	}

	return true,nil
}


//查询用户文件元数据列表
func QueryUserFileMeta(username string, limitNum int)  ([]UserFileTable,error){
	stmt, err := DBConn().Prepare(`SELECT user_name,file_name,file_sha1,file_size,
			create_at,update_at,file_addr from user_file where user_name=? limit ?`)
	if err != nil {
		log.Printf("查询用户%s文件信息失败:%s\n",username,err.Error())
		return  nil, err
	}
	defer stmt.Close()

	var files []UserFileTable
	rows, err := stmt.Query(username,limitNum)
	if err != nil {
		return  files, nil
	}

	for rows.Next() {
		var file UserFileTable
		var upLoadTime time.Time
		var updateTime time.Time
		err := rows.Scan(&file.UserName,&file.FileName,&file.FileHash,&file.FileSize,
			&upLoadTime,&updateTime,&file.FileAddr)

		//format: 2019-01-14T21:51:18+08:00 -> 2019-01-14 21:51:18
		file.UploadAt = upLoadTime.Format("2006-01-02 15:04:05")
		file.LastUpdated = updateTime.Format("2006-01-02 15:04:05")

		if err != nil {
			break
		}
		files = append(files,file)
	}
	return files, err
}
