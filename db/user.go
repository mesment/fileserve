package db
import (
	"errors"
	"log"
	"time"
)

//用户注册

type User struct {
	UserName string
	PassWord string
	Email 	string
	Phone 	string
	SignUpAt string
	LastActiveAt string
	Status int
}

func UserSignUp(username string, password string) (bool,error) {

	stmt, err := DBConn().Prepare(`insert  into 
								user (user_name,user_passwd) values (?,?)`)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return false,err
	}
	res, err := stmt.Exec(username,password)
	if err != nil {
		log.Println(err)
		return false,err
	}

	rowsAffect,err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return false,err
	}
	if rowsAffect <= 0 {
		return false, errors.New("该用户已存在！")
	}
	log.Println("新增用户成功")
	return true, nil
}


//用户登录
func UserSignIn(username, passwd string) (bool, error) {
	stmt, err := DBConn().Prepare(`select user_passwd from user 
	where user_name = ? limit 1`)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return false, err
	}
	var userPasswd string

	err = stmt.QueryRow(username).Scan(&userPasswd)
	if err != nil {
		log.Println(err)
		return false, err
	}
	// 对比密码
	if passwd != userPasswd {
		return false, errors.New("用户密码不正确")
	}
	return true, nil

}

//用户是否存在
func IsUserExist(username string,password string) bool {
	stmt, err := DBConn().Prepare(`select count(*) from user 
	where user_name = ? and user_passwd = ? limit 1`)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return false
	}
	var count int
	err = stmt.QueryRow(username,password).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}
	if count <= 0  {
		return false
	}
		return true
}

func IsUserNameExist(username string) bool {
	stmt, err := DBConn().Prepare(`select count(*) from user 
	where user_name = ? limit 1`)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return false
	}
	var count int
	err = stmt.QueryRow(username).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}
	if count <= 0  {
		return false
	}
	return true
}


//获取用户信息
func GetUserInfo(username string) ( User,  error) {
	stmt, err := DBConn().Prepare(`select user_name,signup_at,last_active from user 
	where user_name = ? limit 1`)
	user := User{}
	if err != nil {
		log.Println(err.Error())
		return user, err
	}
	defer stmt.Close()


	//改变时间显示格式：2019-01-14T21:51:18+08:00 -> 2019-01-14 21:51:18
	var signTime time.Time
	var lastActive time.Time
	err = stmt.QueryRow(username).Scan(&user.UserName,&signTime,&lastActive)
	user.SignUpAt = signTime.Format("2006-01-02 13:04:05")
	user.LastActiveAt = lastActive.Format("2006-01-02 13:04:05")
	if err != nil{
		log.Println(err.Error())
		return user ,err
	}
	return user ,nil
}
