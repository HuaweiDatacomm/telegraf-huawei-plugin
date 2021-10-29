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
	tEscape     = "\t"
	quoteEscape = "\""
	nEscape     = "\n"
	midEscape   = "-"
	downEscape  = "_"
	slashEscape = "/"
)

func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		check := strings.HasSuffix(dirPth, ".proto")
		if check {
			files = append(files, dirPth)
			return files, nil, err
		} else {
			return files, nil, err
		}

	}
	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetFilesAndDirs(dirPth + PthSep + fi.Name())
		} else {

			ok := strings.HasSuffix(fi.Name(), ".proto")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}
	return files, dirs, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func main() {

	argsWithProg := os.Args

	if len(argsWithProg) != 4 {
		fmt.Println("not correct args")
		return
	}
	protoFilePath := os.Args[1]
	pbgoFilePath := os.Args[2]
	huaweiTelemetryPath := os.Args[3]
	files, _, _ := GetFilesAndDirs(protoFilePath)

	_, errProto := os.Stat(pbgoFilePath)
	if errProto != nil {
		fmt.Println("not find go file path")
		return
	}

	if files == nil {
		fmt.Println("no proto file find !")

	}

	TelemetryExist := Exists(huaweiTelemetryPath)
	if TelemetryExist == false {
		fmt.Println("not find HuaweiTelemetry.go!")
		return
	}

	for _, dir := range files {
		var goFileName = huaweiTelemetryPath
		var dirName = dir

		dirNew := strings.Replace(dirName, midEscape, downEscape, -1)
		dirNewList := strings.Split(dirNew, slashEscape)
		dirFin := strings.Split(dirNewList[len(dirNewList)-1], ".")[0]

		dirPaths, _ := filepath.Split(dirName)
		cmd := exec.Command("protoc", "--go_out=plugins=grpc:.", "--proto_path="+dirPaths, dirName)
		fmt.Println(cmd)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}

		dirCurrent, err := filepath.Abs(filepath.Dir(os.Args[0]))
		dirPbGOName := strings.Replace(dirFin, downEscape, midEscape, -1)
		dirOrign := dirCurrent + "/" + dirPbGOName + ".pb.go"
		mkdir_path := pbgoFilePath + "/" + dirFin
		cmd_mkdir := exec.Command("mkdir", mkdir_path)
		cmd_mkdir.Run()
		pbMove := mkdir_path + "/" + dirFin + ".pb.go"
		cmd_move := exec.Command("mv", dirOrign, pbMove)
		errMove := cmd_move.Run()
		if errMove != nil {
			fmt.Println("Error: can't move go file to specify dir,please check filepath and proto file")
			return
		}

		_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		fmt.Printf("warning: %s"+nEscape, errStr)

		telemetry_replace(dirFin, goFileName)
	}
}

func telemetry_replace(dirNew, goFileName string) {
	file, err := os.Open(goFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	dirEndList := strings.Split(dirNew, downEscape)
	dirEnd := dirEndList[len(dirEndList)-1]
	newDirEnd := strings.ToUpper(dirEnd[:1]) + dirEnd[1:]
	for scanner.Scan() {
		var lineText = scanner.Text()
		lines = append(lines, lineText)
		if strings.Contains(lineText, "import (") {
			var impName = tEscape + quoteEscape + "github.com/influxdata/telegraf/plugins/common/telemetry_proto/" + dirNew + quoteEscape
			lines = append(lines, impName)
		} else if strings.Contains(lineText, "PathKey{ProtoPath: \"huawei_debug.Debug\", Version: \"1.0\"}: []reflect.Type{reflect.TypeOf((*huawei_debug.Debug)(nil))}") {
			var var_name = tEscape + "PathKey{ProtoPath: " + quoteEscape + dirNew + "." + newDirEnd + quoteEscape + ", Version: " + quoteEscape + "1.0" + quoteEscape + "}: []reflect.Type{reflect.TypeOf((*" + dirNew + "." + newDirEnd + ")(nil))},"
			lines = append(lines, var_name)
		}
	}
	writeGoFile(lines, goFileName)
}

func writeGoFile(lines []string, goFileName string) {

	newArr := RemoveRepeatedElement(lines)
	strText := strings.Join(newArr, nEscape)
	content := []byte(strText)
	err := ioutil.WriteFile(goFileName, content, 0644)
	if err != nil {
		panic(err)
	}
}

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
