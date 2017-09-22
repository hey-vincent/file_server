package cloud

import (
	"io/ioutil"
	"fmt"
	"os"
	"net/http"
	"html/template"
	"encoding/json"
)

// 文件结构
type FILE_INFO struct {
	Name  string
	Src   string
	Dir   string
	Size  int64
	Mode  int8
	IsDir bool
}

func (elm *FILE_INFO) construct(info os.FileInfo) {
	if nil == info {
		return
	}
	elm.Name = info.Name()
	elm.Size = info.Size()
	elm.IsDir = info.IsDir()
	elm.Mode = int8(info.Mode())
}

// List all the files inside the giving directory
func ListFiles(uri string) []os.FileInfo {
	fn, err := ioutil.ReadDir(uri)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	op := make([]os.FileInfo, 0, len(fn))
	for _, info := range fn {
		if info.Name()[0] != '.' {
			op = append(op, info)
		}
	}

	return op

}

/**

 */
func getFileList(uri string) []FILE_INFO {
	fn, err := ioutil.ReadDir(uri)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fileInfoArray := make([]FILE_INFO, 0, len(fn))
	for _, one := range fn {
		elm := &FILE_INFO{}
		elm.construct(one)
		elm.Dir = uri
		fileInfoArray = append(fileInfoArray, *elm)
	}

	return fileInfoArray
}

// Read file content
func GetFileContent(uri string) []byte {

	ct, err := ioutil.ReadFile(uri)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return ct
}

// file reviewing server
func Serve(w http.ResponseWriter, path string) {
	file_nmae := root_dir + path
	fmt.Println("路径：", file_nmae)

	fi, err := os.Stat(file_nmae)
	if err != nil {
		fmt.Println("RETURN", err.Error())
		return
	}

	// 如果是文件
	if !fi.IsDir() {

		// 文件大小大于3M
		if fi.Size() > 1024*1024*3 {
			w.Write([]byte("Sorry. File si too large to read."))
			return
		}

		cont := GetFileContent(file_nmae)
		if nil != cont {
			w.Write(cont)
			fmt.Println("RETURN", 2)
			return
		}
		fmt.Println("RETURN", 3)

		return
	}

	fList := make([]FILE_INFO, 0, 10)

	for _, v := range SortContent(ListFiles(root_dir + path)) {
		uri := path + v.Name()
		if v.IsDir() {
			uri += "/"
		}
		fList = append(fList, FILE_INFO{Name: v.Name(), Src: uri, IsDir: v.IsDir()})
	}
	//fmt.Println("内容",fList)

	tpl, err := template.ParseFiles("./static/template/list.html")
	if nil != err {
		fmt.Println(err.Error())
		return
	}

	tpl.Execute(w, fList)

}

/**
 * 将目录排序,文件夹在前 文件在后
 */
func SortContent(flist []os.FileInfo) []os.FileInfo {
	for key, val := range flist {
		if !val.IsDir() {
			for x := key + 1; x < len(flist); x++ {
				if flist[x].IsDir() {
					flist[key], flist[x] = flist[x], flist[key]
					break
				}
				if x == len(flist)-1 {
					return flist
				}
			}

		}
	}
	return flist
}

// reflection
func ResortFiles(fileList []FILE_INFO) []FILE_INFO {
	for key, val := range fileList {
		if !val.IsDir {
			for x := key + 1; x < len(fileList); x++ {
				if fileList[x].IsDir {
					fileList[key], fileList[x] = fileList[x], fileList[key]
					break
				}
				if x == len(fileList)-1 {
					return fileList
				}
			}

		}
	}
	return fileList
}

/**
	取一个给定路径下的文件列表
 */
func GetFilesListOf(path string) []byte {

	dir := root_dir + path
	fList := getFileList(dir)

	fList = ResortFiles(fList)
	flJson, _ := json.Marshal(fList)
	fmt.Println(flJson)
	return flJson
}
