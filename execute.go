package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	tEscape = "\t"
	quoteEscape = "\""
	nEscape = "\n"
	midEscape = "-"
	downEscape = "_"
	slashEscape = "/"
)

func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		check := strings.HasSuffix(dirPth, ".proto")
		if check {
			files = append(files,dirPth)
			return files, nil, err
		} else {
			return files, nil,err
		}

	}
	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetFilesAndDirs(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".proto")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}
	return files, dirs, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func main() {
	//获取输入的参数，第一个是proto文件所在目录，第二个应该是生成的go脚本的路径，第三个是HUAWEITelemetry.go文件所在目录
	argsWithProg := os.Args
	//入参数量判断
	if(len(argsWithProg) != 4) {
		fmt.Println("not correct args")
		return
	}
	protoFilePath := os.Args[1]
	pbgoFilePath := os.Args[2]
	huaweiTelemetryPath := os.Args[3]
	files, _, _ := GetFilesAndDirs(protoFilePath)
	//判断生成的路径是否存在
	_, errProto := os.Stat(pbgoFilePath)    //os.Stat获取文件信息
	if errProto != nil {
		fmt.Println("not find go file path")
		return
	}

	if files == nil{
		fmt.Println("no proto file find !")

	}

	//判断HuaweiTelemetroy路径是否存在
	TelemetryExist := Exists(huaweiTelemetryPath)
	if TelemetryExist == false{
		fmt.Println("not find HuaweiTelemetry.go!")
		return
	}

	for _, dir := range files{
		var goFileName = huaweiTelemetryPath
		var dirName = dir
		//处理下划线
		dirNew := strings.Replace(dirName, midEscape, downEscape, -1)
		dirNewList := strings.Split(dirNew,slashEscape)
		dirFin := strings.Split(dirNewList[len(dirNewList)-1], ".")[0]
		//获取proto文件的路径，作为proto_path的参数，主要是为了防止报错
		dirPaths,_ := filepath.Split(dirName)
		cmd := exec.Command("protoc","--go_out=plugins=grpc:.","--proto_path="+dirPaths,dirName)
		fmt.Println(cmd)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		//移动生成的pb.go到指定目录下
		dirCurrent, err := filepath.Abs(filepath.Dir(os.Args[0]))
		dirPbGOName := strings.Replace(dirFin,downEscape,midEscape,-1)
		dirOrign := dirCurrent + "/" + dirPbGOName + ".pb.go"
		mkdir_path := pbgoFilePath+"/"+dirFin
		cmd_mkdir := exec.Command("mkdir",mkdir_path)
		cmd_mkdir.Run()
		pbMove := mkdir_path + "/" + dirFin+".pb.go"
		cmd_move := exec.Command("mv",dirOrign,pbMove)
		errMove := cmd_move.Run()
		if errMove != nil{
			fmt.Println("Error: can't move go file to specify dir,please check filepath and proto file")
			return
		}

		_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		fmt.Printf("warning: %s"+nEscape, errStr)

		//HUAWEITelemetry.go文件内容替换
		telemetry_replace(dirFin, goFileName)}
}

func telemetry_replace(dirNew,goFileName string){
	file ,err := os.Open(goFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	//是否有下一行
	var lines [] string
	//没有下划线的情况
	dirEndList := strings.Split(dirNew,midEscape)
	dirEnd := dirEndList[len(dirEndList)-1]
	for scanner.Scan() {
		var lineText = scanner.Text()
		lines = append(lines, lineText)
		if strings.Contains(lineText, "import (") {
			var impName = tEscape+quoteEscape+"github.com/influxdata/telegraf/plugins/common/telemetry_proto/" + dirNew + quoteEscape
			lines = append(lines, impName)
		}else if strings.Contains(lineText,"PathKey{ProtoPath: \"huawei_debug.Debug\", Version: \"1.0\"}: []reflect.Type{reflect.TypeOf((*huawei_debug.Debug)(nil))}"){
			var var_name = tEscape+"PathKey{ProtoPath: "+quoteEscape+ dirNew +"."+ dirEnd +quoteEscape+", Version: "+quoteEscape+"1.0"+quoteEscape+"}: []reflect.Type{reflect.TypeOf((*"+ dirNew +"."+ dirEnd +")(nil))},"
			lines = append(lines, var_name)
		}
	}
	writeGoFile(lines,goFileName)
}

func writeGoFile(lines []string,goFileName string) {
	//对lines进行去重
	newArr := RemoveRepeatedElement(lines)
	strText := strings.Join(newArr, nEscape)
	content := []byte(strText)
	err := ioutil.WriteFile(goFileName, content, 0644)
	if err != nil {
		panic(err)
	}
}

//列表去重
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}