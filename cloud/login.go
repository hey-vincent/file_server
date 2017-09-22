package cloud

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	authority_list = map[string]string{
		"wensheng": "wensheng123",
		"vincent":  "helloworld", }
	// md5 签名
	key = "YOU1SHOULD2CONTROL3YOURSELF"
	
	cookie_max_age = 60 * 10
)


func IsLoggedin(r *http.Request) bool{
	fmt.Println("请求信息" , r)
	
	
	user_name,err1 := r.Cookie("username")
	user_pass,err2 := r.Cookie("password")
	
	
	if nil != err1 {
		fmt.Println(err1.Error())
		return false
	}
	
	if nil != err2{
		fmt.Println(err2.Error())
		return false
	}
	
	
	
	if user_name.Value == "wensheng"  && user_pass.Value == "wensheng123"{
		return true
	}
	
	return false
}

// 显示登录页
func ShowLoginPage(w http.ResponseWriter){
	p, _ := filepath.Abs("./")
	fmt.Println("当前root", p)
	tpl, err := template.ParseFiles( "./static/template/login.html")
	//fmt.Println(tpl)
	if nil != err{
		//fmt.Println(err.Error())
		return
	}
	tpl.Execute(w, nil)
}


// 检查用户名密码是否匹配
func CheckAuth(w http.ResponseWriter, user_name , password string)bool{

	pass , ok := authority_list[user_name]
	if !ok || pass != password{
		return false
	}

	// 登录成功
	http.SetCookie(w , &http.Cookie{Name:"password",Value:password ,  MaxAge: cookie_max_age})
	http.SetCookie(w , &http.Cookie{Name:"username",Value:user_name , MaxAge:cookie_max_age})
	return true
}



