package main

import (
	"github.com/mesment/fileserver/pkg/setting"
	"github.com/mesment/fileserver/handler"
	"github.com/mesment/fileserver/db"
	"net/http"
)

func InitRouter()  {
	//配置静态资源
	http.Handle("/static/",http.StripPrefix("/static/",http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload",handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/success",handler.UploadSuccessHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadFileHandler)
	http.HandleFunc("/file/update",handler.UpdateFileMetaHandler)
	http.HandleFunc("/file/delete",handler.DeleteFileHandler)
	http.HandleFunc("/file/query",handler.HTTPInterceptor(handler.QueryUserFileHandler))
	http.HandleFunc("/user/signup",handler.UserSignUpHandler)
	http.HandleFunc("/user/signin",handler.UserSignInHandler)
	http.HandleFunc("/user/info",handler.HTTPInterceptor(handler.UserInfoHandler))

	http.HandleFunc("/file/fastupload",handler.HTTPInterceptor(handler.TryFastUpLoadHandler))



}

func main() {
	setting.Setup()
	handler.Setup()
	db.Setup()

	
	InitRouter()
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic("启动服务失败")
	}
}
