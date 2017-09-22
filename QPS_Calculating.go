package main

import (
	"regexp"
	"fmt"
	"time"
	"net/http"
	"sync"
)

func main(){
	strongWeb()
	//grep()
	
	
}

func grep(){
	for {
		time.Sleep(1 * time.Second)
		str := ""
		fmt.Println("input source string for matching ^xy/:")
		fmt.Scanf("%s", &str)
		yes, _ := regexp.Match("^/[45][0-9][0-9].html", []byte(str))
		
		fmt.Println(yes)
		
	}
}

var ch = make(chan int  , 4000)
var chStop = make(chan int , 1)		// routine 退出通知
var time_start time.Time		//开始时间
var request_url	= "http://127.0.0.1:8000/login"
var query_control = 2000
var query_count = 0

var mutex = sync.Mutex{}


func request_routine(){

	mutex.Lock()
	if time_start.IsZero(){
		time_start = time.Now()
		fmt.Println("设置时间",time_start.String())
	}
	mutex.Unlock()

	for{
		_, err := http.Get(request_url)

		if nil != err  {
			fmt.Println("Request Error：" , err.Error())
			break
		}

		mutex.Lock()
		if len(chStop) > 0 {
			fmt.Println("Quit goroutine")
			break
		}

		query_count += 1
		mutex.Unlock()
		ch <-1
	}

}

func strongWeb(){

	for i:=0;i<100;i++{
		go request_routine()
	}

	stop := false

	for{

		if stop {
			break
		}

		//fmt.Println(len(ch))
		if len(ch) >= query_control {
			chStop <- 1

			calQps(time_start, time.Now() , uint64(query_count) )

			stop = true
			break
		}
		time.Sleep(1)

	}

}

// calculating QPS
func calQps(time_start , time_end time.Time , queries uint64) float64{
	duration := float64(time_end.UnixNano() - time_start.UnixNano()) / 1e9
	fmt.Println("耗时：", duration , "请求：" , queries)
	qps := float64(queries) / duration
	fmt.Println("QPS: ",qps)
	return qps
}