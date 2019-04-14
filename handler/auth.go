package handler

import (
	"log"
	"github.com/gin-gonic/gin"
	dblayer "github.com/mesment/fileserver/db"
	"github.com/astaxie/beego/validation"
	"github.com/mesment/fileserver/utils"
	"net/http"
	pkge "github.com/mesment/fileserver/pkg/errors"
)


//HTTPInterceptor : http请求拦截器
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter,r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username) < 3 || !utils.IsTokenValid(username,token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h(w,r)
		})
}


type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := pkge.INVALID_PARAMS
	if ok {
		isExist := CheckUserExist(username, password)
		if isExist {
			token, err := utils.GenerateToken(username, password)
			if err != nil {
				code = pkge.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token

				code = pkge.SUCCESS
			}

		} else {
			code = pkge.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : pkge.GetMsg(code),
		"data" : data,
	})
}





func CheckUserExist(username, password string) bool {

	if exist := dblayer.IsUserExist(username,password);exist {
		return true
	}

	return false
}


