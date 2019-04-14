package db

import (
	"database/sql"
	"log"
	"os"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type FileTable struct {
	FileSha1	sql.NullString
	FileMD5		sql.NullString
	FileName 	sql.NullString
	FileAddr	sql.NullString
	FileSize 	sql.NullInt64
}


//OnFileUploadFinish：保存文件元数据到数据库
func OnFileUploadFinish(sha1,md5, filename string, filesize int64, fileaddr string ) bool {
	stmt, err := DBConn().Prepare(`insert ignore into filemeta(file_sha1,file_md5,file_name,
										file_size,file_addr) values (?,?,?,?,?)`)
	if err != nil {
		log.Printf("插入文件信息失败：%s\n",err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(sha1,md5,filename,filesize,fileaddr)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	ret.LastInsertId()
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			log.Printf("影响数目:%d\n",rf)
			log.Printf("文件SHA1:%s已存在\n",sha1)
		}
		return true
	}
	return false
}

//UpdateFileMeta：更新数据库中文件元数据
func UpdateFileMeta(sha1 string,name string, addr string,size int64) bool {
	stmt, err := DBConn().Prepare(`update filemeta set file_name=?,
						file_addr=?,file_size=? where file_sha1=?` )
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_,err = stmt.Exec(name,addr,size,sha1)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}



//GetFileMeta:从数据库查询文件元数据
func GetFileMeta(sha1 string) (*FileTable, error){
	stmt,err := DBConn().Prepare(`select file_sha1,file_md5, file_name, file_addr, file_size
	  							from filemeta where file_sha1 = ?  and status='0' limit 1`)
	if err != nil {
		log.Println("查询文件元数据失败：",err.Error())
		return nil, err
	}
	defer  stmt.Close()

	ft := FileTable{}
	err = stmt.QueryRow(sha1).Scan(&ft.FileSha1,&ft.FileMD5,&ft.FileName,&ft.FileAddr,&ft.FileSize)
	if err != nil {
		log.Println(err.Error())
		return nil ,err
	}

	return &ft, nil
}

//DeleteFileMetaDB：从数据库中删除文件元数据
func DeleteFileMetaDB(sha1 string) bool {
	stmt, err := DBConn().Prepare(`delete from filemeta where file_sha1 = ?`)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_,err = stmt.Exec(sha1)
	if err != nil {
		return false
	}
	return true
}
