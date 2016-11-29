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
const rootPath = "//gfs/Shares/00_全社共有/"

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
	file, _ := os.OpenFile(currentDir+"/capacity.csv", os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	writer := bufio.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))

	var csv string
	for k, v := range capacityMap {
		content := rootPath + k + "," + strconv.FormatInt(v.Size, 10) + "," + strconv.Itoa(v.Count) + "\n"
		writer.WriteString(content)
		csv += content
	}

	fmt.Println(csv)
	writer.Flush()
}

// 階層探索します
func searchRoot() {
	files, _ := ioutil.ReadDir(rootPath)
	for _, file := range files {
		if file.IsDir() {
			fileName := file.Name()
			fmt.Println(fileName)
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
			fmt.Println(dir1 + "/" + fileName)
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
