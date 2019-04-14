package handler

import (
	"encoding/json"
	"fmt"
	dbLayer "github.com/mesment/fileserver/db"
	pkge "github.com/mesment/fileserver/pkg/errors"
	"github.com/mesment/fileserver/utils"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	passwd_salt = "&^%$#*&^%"
)

//用户注册
func UserSignUpHandler(w http.ResponseWriter, r *http.Request)  {
	switch r.Method {
	case http.MethodGet:
		data,err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			log.Println("读取注册页面失败：%s\n",err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	case http.MethodPost:
		r.ParseForm()
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		//check username, password
		if len(username)<=3 || username  == "" || password == "" || len(password) <=3 {
			w.Write([]byte("Fail"))
			return
		}
		fmt.Printf("username:%s,passowrd:%s",username,password)
		cryptpasswd := utils.Sha1([]byte(password+ passwd_salt))

		//查看用户名是否已存在，已存在报错
		if exist := dbLayer.IsUserNameExist(username); exist {

			w.Write([]byte("Fail"))
		}
		//新增用户
		_,err := dbLayer.UserSignUp(username,cryptpasswd)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			errinfo :=fmt.Sprintf("注册失败：%s",err.Error())
			fmt.Println(errinfo)
			w.Write([]byte(errinfo))
			return
		}
		w.Write([]byte("SUCCESS"))

	}

}


//用户登录
func UserSignInHandler(w http.ResponseWriter, r *http.Request)  {
	switch r.Method {
	case http.MethodGet:
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			log.Println("读取注册页面失败：%s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	case http.MethodPost:
		r.ParseForm()
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		//check username, password
		if len(username)<=3 || username  == "" || password == "" || len(password) <=3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := 0
		var token string
		var err error
		cryptpasswd := utils.Sha1([]byte(password+ passwd_salt))
		if pass := CheckUserExist(username,cryptpasswd); pass {

			//生成访问凭证
			token, err = utils.GenerateToken(username, password)

			if err != nil {
				log.Println(err)
				code = pkge.ERROR_AUTH_TOKEN
			} else {
				code = pkge.SUCCESS
			}


		} else {
			//用户不存在，无权限
			w.WriteHeader(http.StatusForbidden)
			return
		}
		resp :=utils.RespMsg{
			Code: code,
			Msg:pkge.GetMsg(code),
			Data: struct {
				Location string
				Username string
				Token string
			}{
				Location:"http://"+ r.Host + "/static/view/home.html",
				Username:username,
				Token: token,
			},
		}
		fmt.Printf("Resp:%s\n",resp.JSONBytes())
		w.Write(resp.JSONBytes())

	}
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request)  {
	//1、解析请求参数
	r.ParseForm()
	token := r.Form.Get("token")
	username := r.Form.Get("username")
	//2. 验证token
	if !utils.IsTokenValid(username,token) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//查询用户信息
	log.Printf(" recv: username:%s,token:%s\n",username,token)
	user, err := dbLayer.GetUserInfo(username)
	log.Printf("username:%s,signupat:%s",user.UserName,user.SignUpAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//4、组装并响应用户数据
	resp:= utils.RespMsg{
		Code:0,
		Msg:"OK",
		Data:user,
	}

	w.Write(resp.JSONBytes())
}


//查询用户的文件
func QueryUserFileHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	if valid := utils.IsTokenValid(username,token); !valid {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	userfile,err := dbLayer.QueryUserFileMeta(username,15)
	if err != nil {
		log.Println(err)
	}

	data, err := json.Marshal(userfile)

	w.Write(data)
}

