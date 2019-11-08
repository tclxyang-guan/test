package utils

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/kataras/iris/v12"
	"io"
	"net/http"
	"os"
)

type FileUtil struct {
}

func NewFileUtil() *FileUtil {
	return &FileUtil{}
}

//文件上传
func (f *FileUtil) UploadFile(ctx iris.Context, path string) (string, error) {
	//获取文件内容 要这样获取
	file, head, err := ctx.FormFile("file")
	filename := head.Filename
	if err != nil {
		return "", err
	}
	defer file.Close()
	//创建文件
	fW, err := os.Create(path + head.Filename)
	if err != nil {
		fmt.Println("文件创建失败")
		return "", err
	}
	defer fW.Close()
	_, err = io.Copy(fW, file)
	if err != nil {
		fmt.Println("文件保存失败")
		return "", err
	}
	return filename, nil
}

//文件下载
func (f *FileUtil) DownLoadFile(ctx iris.Context, allFileName string, newFileName string) error {
	return ctx.SendFile(allFileName, newFileName)
}

//读取Excel
func (f *FileUtil) ReadExcel(path string, fileName string) ([][]string, error) {
	var data [][]string
	xlsx, err := excelize.OpenFile(path + fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Get value from cell by given worksheet name and axis.
	// Get all the rows in the Sheet1.
	rows := xlsx.GetRows("Sheet1")
	for _, row := range rows {
		var rdata []string
		for _, colCell := range row {
			rdata = append(rdata, colCell)
		}
		data = append(data, rdata)
	}
	return data, nil
}

//多文件上传 路径 关联的id 反馈文件还是核实文件
func (f *FileUtil) UploadManyFile(r *http.Request, path string) {
	//设置内存大小
	r.ParseMultipartForm(32 << 20)
	//获取上传的文件组
	files := r.MultipartForm.File["file"]
	len := len(files)
	for i := 0; i < len; i++ {
		//打开上传文件
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Println("打开失败")
		}
		//创建上传目录
		os.Mkdir(path, os.ModePerm)
		//创建上传文件
		cur, err := os.Create(path + "文件名")
		defer cur.Close()
		if err != nil {
			fmt.Println("创建失败")
		}
		io.Copy(cur, file)
	}
	return
}
