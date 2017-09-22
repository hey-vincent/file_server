package main

import (
	"runtime"
	"net/http"
	"net"
	"MyCloud/cloud"
	"regexp"
	"fmt"
	"net/url"
	"html/template"
	"time"
)

const domain = ":8000"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	StartSvr()
	
}

// 开启服务
func StartSvr() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	
	http.HandleFunc("/", Route)
	lst, _ := net.Listen("tcp", domain)
	svr := http.Server{Handler: nil}
	svr.SetKeepAlivesEnabled(false)
	svr.Serve(lst)
	
}

//服务路由
func Route(w http.ResponseWriter, r *http.Request) {
	
	
	switch {
	case r.URL.Path == "/login":
		if r.Method == "POST" {
			r.ParseForm()
			form := r.PostForm
			// 验证失败
			if ! cloud.CheckAuth(w, form.Get("username"), form.Get("password")) {
				return
			}
			// 验证成功
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			fmt.Println(r.Method , r.Proto , r.Header.Get("Accept"))
			// 未登录显示登录页面
			if !cloud.IsLoggedin(r) {
				cloud.ShowLoginPage(w)
				return
			}
			// 如果已经登录 跳转到首页
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusSeeOther)
		}
		break
	
	case  matchURL("^/link/" , r.URL.Path):
		if r.Method == "GET"{
			defer func(){
				fmt.Println("defer occured")
				time.Sleep(time.Second * 3)
				w.Header().Set("Location" , fmt.Sprintf("/%d.html" ,http.StatusMethodNotAllowed ))
				w.WriteHeader(http.StatusSeeOther)
			}()
			
			referer ,_err := url.Parse( r.Referer() )
			if _err != nil || referer == nil{
				fmt.Println("Error while getting referer")
			}
		
			if referer.Hostname() != "localhost"{
				return
			}
			
			// parse query string
			mapQuery,err :=url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			
			token := mapQuery.Get("token")
			fmt.Println(token)
			if  string(token) == ""{
				fmt.Println("Missing token")
			}
			
			path := r.URL.Path[len("/link") :]
			cloud.Serve(w, path)
		}else{

		}

		break
	
	case matchURL("^/[45][0-9][0-9].html" , r.URL.Path):
		fmt.Println("ERROR PAGE")
		tpl, err := template.ParseFiles("./static/template/4xx.html")
		if nil != err{
			fmt.Println(err.Error())
			return
		}
		
		tpl.Execute(w, w.Header().Get("status"))
		break


	case matchURL("^/file/", r.URL.Path):
		if r.Method != "POST"{
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// get files under the given path
		path := r.PostForm.Get("path")
		jsn := cloud.GetFilesListOf(path)
		w.Write([]byte(string(jsn)))
		break

	default:
		//if !cloud.IsLoggedin(r) {
		//
		//	//w.Header().Set("Location", "/login")
		//	//w.WriteHeader(http.StatusSeeOther)
		//	//return
		//} else {
			cloud.Serve(w, r.URL.Path)
		//}
		break
	}
	
}

func matchURL(pattern , url string) bool{
	yes , err := regexp.Match(pattern , []byte(url))
	if nil != err {
		return false
	}
	return yes
}
