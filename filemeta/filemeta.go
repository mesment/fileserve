package filemeta

import  (
	mydb "github.com/mesment/fileserver/db"
)

type FileMeta struct {
	FileSha1		string      //文件sha1
	FileMD5 		string		//文件md5
	FileName		string		//文件名
	FileLocation 	string		//文件存储地址
	FileSize		int64		//文件大小
	CreateTime 		string		//创建时间
}

// Map存储用户上传文件的元数据
var fileMetaStore map[string]FileMeta;

func init()  {
	fileMetaStore = make(map[string]FileMeta)
}

//GetFileMeta:根据文件sha1查找文件元数据
func GetFileMeta(sha string) *FileMeta {
	meta,ok := fileMetaStore[sha]
	if !ok {
		return nil
	}
	return &meta
}


//AddFileMeta：添加文件元数据
func AddFileMeta( meta *FileMeta)  {
	fileMetaStore[meta.FileSha1] = *meta
}

//UpdateMeta：更新文件元数据
func UpdateFileMeta( newMeta *FileMeta)  {
	fileMetaStore[newMeta.FileSha1] = *newMeta
}

//DeleteFileMeta：删除文件元数据
func DeleteFileMeta(sha1 string)  {
	delete(fileMetaStore, sha1)
}


//GetFileMetaDB:从数据库查找文件元数据
func GetFileMetaDB(sha1 string) (*FileMeta,error)  {
	filetab, err := mydb.GetFileMeta(sha1)
	if err != nil {
		return nil,err
	}
	meta := FileMeta{
		FileSha1:filetab.FileSha1.String,
		FileMD5:filetab.FileMD5.String,
		FileName:filetab.FileName.String,
		FileLocation:filetab.FileAddr.String,
		FileSize:filetab.FileSize.Int64,
	}
	return &meta,nil
}

//AddFileMetaDB：往数据库添加文件元数据
func AddFileMetaDB ( meta *FileMeta) bool {
	return mydb.OnFileUploadFinish(meta.FileSha1,meta.FileMD5,meta.FileName,meta.FileSize,meta.FileLocation)
}

//UpdateMeta：更新数据库的文件元数据
func UpdateFileMetaDB( newMeta *FileMeta) bool {
	return mydb.UpdateFileMeta(newMeta.FileSha1,newMeta.FileName,newMeta.FileLocation,newMeta.FileSize)
}

//DeleteFileMetaDB:从数据库删除文件元数据
func DeleteFileMetaDB(sha1 string) bool {
	return mydb.DeleteFileMetaDB(sha1)
}

