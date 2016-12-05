package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// ルートディレクトリパス
const rootPath = "\\\\gfs\\Shares\\00_全社共有\\"

//const rootPath = "C:/Workspace/Es/svn/"

// capacity 階層の容量を表す
type capacity struct {
	Size  int64
	Count int
}

// 階層容量マップ
var capacityMap = map[string]capacity{}

// エントリポイント
func main() {
	start := time.Now()

	searchRoot()
	output()

	end := time.Now()
	fmt.Printf("%f(s)", (end.Sub(start)).Seconds())
}

// csvを出力します
func output() {
	currentDir, _ := os.Getwd()
	file, _ := os.Create(currentDir + "/capacity.csv")
	defer file.Close()
	//writer := bufio.NewWriterSize(file, 20000)
	writer := bufio.NewWriterSize(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()), 20000)
	for k, v := range capacityMap {
		writer.WriteString(rootPath + k + "," + strconv.FormatInt(v.Size, 10) + "," + strconv.Itoa(v.Count) + "\n")
	}
	writer.Flush()
}

func test() {
	currentDir, _ := os.Getwd()
	txt, _ := os.Open(currentDir + "/hoge.txt")
	defer txt.Close()
	scanner := bufio.NewScanner(txt)

	file, _ := os.Create(currentDir + "/hoge.csv")
	defer file.Close()
	writer := bufio.NewWriterSize(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()), 20000)

	for scanner.Scan() {
		content := scanner.Text()
		bytes, _ := writer.WriteString(content + "\n")
		fmt.Println(bytes)
	}
	writer.Flush()
}

// 階層探索します
func searchRoot() {
	files, _ := ioutil.ReadDir(rootPath)
	for _, file := range files {
		if file.IsDir() {
			fileName := file.Name()
			capacityMap[fileName] = capacity{Size: 0, Count: 0}
			searchFirst(fileName)
		}
	}
}

func searchFirst(dir1 string) {
	files, _ := ioutil.ReadDir(rootPath + "/" + dir1)
	for _, file := range files {
		if file.IsDir() {
			fileName := file.Name()
			capacityMap[dir1+"/"+fileName] = capacity{Size: 0, Count: 0}
			searchSecond(dir1, fileName)
		} else {
			updateMap(dir1, file.Size())
		}
	}
}

func searchSecond(dir1 string, dir2 string) {
	secondPath := rootPath + "/" + dir1 + "/" + dir2
	files, _ := ioutil.ReadDir(secondPath)
	for _, file := range files {
		if file.IsDir() {
			searchOther(dir1, dir2, secondPath+"/"+file.Name())
		} else {
			fileSize := file.Size()
			updateMap(dir1, fileSize)
			updateMap(dir1+"/"+dir2, fileSize)
		}
	}
}

func searchOther(dir1 string, dir2 string, targetPath string) {
	files, _ := ioutil.ReadDir(targetPath)
	for _, file := range files {
		if file.IsDir() {
			searchOther(dir1, dir2, targetPath+"/"+file.Name())
		} else {
			fileSize := file.Size()
			updateMap(dir1, fileSize)
			updateMap(dir1+"/"+dir2, fileSize)
		}
	}
}

func updateMap(key string, size int64) {
	value, _ := capacityMap[key]
	value.Count++
	value.Size += size
	capacityMap[key] = value
}
