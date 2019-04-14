package handler

import (
	"encoding/json"
	dblayer "github.com/mesment/fileserver/db"
	"github.com/mesment/fileserver/filemeta"
	"github.com/mesment/fileserver/pkg/setting"
	"github.com/mesment/fileserver/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)
//文件存放路径
var FileStorePath string

func Setup()  {
	 FileStorePath = setting.ServerCfg.FileStorePath

}

//UploadHandler:处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		//upload
		data,err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			log.Printf("服务出错:%v",err)
			io.WriteString(w, err.Error())
			return
		}

		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		r.ParseForm()
		//receive file data stream and store
		file, head, err := r.FormFile("file")
		if err != nil {
			log.Printf("接收文件失败")
			return;
		}
		defer file.Close()

		username := r.Form.Get("username")

		//设置文件保存路径
		location := FileStorePath + "/" + head.Filename
		meta :=filemeta.FileMeta{
			FileName:head.Filename,
			FileLocation: location,
		}

		//创建文件句柄来接收文件流
		fd, err := os.Create(location)
		if err != nil {
			log.Printf("创建文件失败%s",err.Error())
			return
		}
		defer fd.Close()

		//拷贝文件到本地
		meta.FileSize, err = io.Copy(fd, file)
		if err != nil {
			log.Printf("保存文件失败%s",err.Error())
			return
		}
		//设置文件元数据
		fd.Seek(0,0)
		meta.FileSha1 = utils.FileSha1(fd)
		fd.Seek(0,0)
		meta.FileMD5 = utils.FileMD5(fd)
		meta.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		//保存文件元数据到数据库
		filemeta.AddFileMetaDB(&meta)

		//保存文献信息到用户信息表
		_,err = dblayer.SaveUserFileInfo(username,meta.FileName,meta.FileSha1,
			meta.FileMD5,meta.FileLocation,meta.FileSize)
		if err != nil {
			log.Println(err)
			data := utils.RespMsg{
				Code:1,
				Msg:"上传文件失败" + err.Error(),
			}
			w.Write(data.JSONBytes())
			return
		}

		http.Redirect(w,r,"/file/upload/success", http.StatusFound)
	}
	
}

//UploadSuccessHandler:文件上传成功
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "文件上传成功")
}


//GetFileMetaHandler:获取文件元数据
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	md5info := r.Form.Get("filehash")
	meta,err := filemeta.GetFileMetaDB(md5info)
	if meta == nil {
		w.WriteHeader(http.StatusNotFound)   //文件信息不存在
		return
	}
	data,err := json.Marshal(meta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}


//DownloadFileHandler:文件下载
func DownloadFileHandler(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	md5info := r.Form.Get("filehash")
	//从数据库查询文件元数据
	meta,err:= filemeta.GetFileMetaDB(md5info)
	if meta == nil {
		w.WriteHeader(http.StatusNotFound)   //文件信息不存在
		return
	}

	fd, err := os.Open(meta.FileLocation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer  fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\"" +meta.FileName+"\"")
	w.Write(data)
}


//UpdateFileMetaHandler:更新文件元数据
func UpdateFileMetaHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	option :=r.Form.Get("op")
	md5info := r.Form.Get("filehash")
	newfilename := r.Form.Get("filename")

	if option != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if (r.Method != "POST") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//从数据库查询文件元数据
	curMeta,err:= filemeta.GetFileMetaDB(md5info)
	if curMeta == nil {
		w.WriteHeader(http.StatusNotFound)   //文件信息不存在
		return
	}
	//更新元数据中文件名
	curMeta.FileName = newfilename
	//更新元数据
	success := filemeta.UpdateFileMetaDB(curMeta)
	if !success {
		log.Printf("更新元数据失败")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(*curMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

//DeleteFileHandler: 删除文件
func DeleteFileHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	md5info := r.Form.Get("filehash")
	//从数据库查询文件元数据
	meta ,err := filemeta.GetFileMetaDB(md5info)
	if meta == nil {
		w.WriteHeader(http.StatusNotFound)   //文件信息不存在
		return
	}
	err = os.Remove(meta.FileLocation)
	if err != nil {
		log.Printf("删除文件失败:%s\n",err.Error())
	}
	//从数据库删除文件元数据
	success := filemeta.DeleteFileMetaDB(md5info)
	if !success {
		log.Printf("从数据库删除元数据失败")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func TryFastUpLoadHandler(w http.ResponseWriter, r *http.Request)  {
	//1 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize := r.Form.Get("filesize")

	fileSize,err := strconv.ParseInt(filesize,10,64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//2 从文件表里查询hash是否存在
	filemeta ,err := filemeta.GetFileMetaDB(filehash)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//3 查不到记录返回失败
	if filemeta == nil{
		resp := utils.RespMsg{
			Code:-1,
			Msg:"秒传失败，请使用普通上传接口",
		}
		w.Write(resp.JSONBytes())
	}

	//4 查到 则将文件信息写入到用户信息表，返回成功
	suc,err := dblayer.SaveUserFileInfo(username,filename,filehash,filemeta.FileMD5,
		filemeta.FileLocation,fileSize)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if suc {
		resp := utils.RespMsg{
			Code:0,
			Msg:"秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	} else {
		resp := utils.RespMsg{
			Code:-2,
			Msg:"秒传失败，稍后再试",
		}
		w.Write(resp.JSONBytes())
	}
}