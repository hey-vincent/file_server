package cloud

import (
	"io/ioutil"
	"fmt"
	"strings"
	"os"
)

var root_dir = ""
var Server_port = ""


func init() {
	mapConf :=  ParseIniFile("./conf/server.ini");
	if nil == mapConf || len(mapConf) < 0{
		fmt.Println("failed to parse ini file")
		os.Exit(1)
		return
	}
	
	dir , ok := mapConf["root_dir"]
	if !ok {
		fmt.Println("root_dir not found")
		os.Exit(1)
		return
	}
	root_dir = dir
	fmt.Println("#Root direction:" , root_dir)
}

func ParseIniFile(file_name string) map[string]string {
	btContent, ok := ioutil.ReadFile(file_name)
	if nil != ok{
		fmt.Println(ok.Error())
		return nil
	}
	
	// to check if 'index' is the tail of current line
	isNewline := func( index int) bool {
		//fmt.Println(string(btContent[index:index+2]))
		if string(btContent[index:index+1]) == "\n" {
			return true
		}
		return false
	}
	
	// to check if 'line' is a comment one
	isComment := func(line string) bool {
		line =  strings.Trim(line, " ")
		
		if len(line) <= 0 {
			return true
		}
		
		switch{
		case len(line) >= 2 && line[0] == '[' && line[len(line) -1 ] == ']':
			return true
		
		case line[:1] == ";":
			return true
		
		case line[:1] == "#":
			return true
		
		default:
			return false
		}
		
		return false
	}
	
	nStart := 0
	var mapConf map[string]string = make(map[string]string , 0)
	
	for  pos,_  := range btContent{
		
		if  isNewline(pos){
			str := btContent[nStart:pos]
			
			if  !isComment(string(str)){
				
				if equal := strings.Index(string(str) , "="); equal > -1 {
					if equal < len(str) -1 {
						mapConf[string(str[:equal])] = string(str[equal+1:])
					}
					
				}
				
			}
			pos += 1
			nStart = pos
		}
	}
	
	
	return mapConf
}