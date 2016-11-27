package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// ルートディレクトリパス
//const rootPath = "//gfs/Shares/00_全社共有/01_セキュリティ対策ソフト/"
const rootPath = "C:/Workspace/golang/"

// 階層を表す
const (
	Root hierarchy = iota
	First
	Other
)

// Hierarchy 階層を表す
type hierarchy int8

// capacity 階層の容量を表す
type capacity struct {
	Size  int64
	Count int
}

// 階層マップ
var hierarchyMap = map[string]capacity{}

// エントリポイント
func main() {
	start := time.Now()

	search(rootPath, "", "")
	output()

	end := time.Now()
	fmt.Printf("%f(s)", (end.Sub(start)).Seconds())
}

// csvを出力します
func output() {
	var csv string
	for k, v := range hierarchyMap {
		csv += rootPath + k + "," + strconv.FormatInt(v.Size, 10) + "," + strconv.Itoa(v.Count) + "\n"
	}

	currentDir, _ := os.Getwd()
	ioutil.WriteFile(currentDir+"/capacity.csv", []byte(csv), os.ModePerm)
}

// 階層探索します
func search(targetPath string, dir1 string, dir2 string) {
	hierarchy := getHierarchy(dir1, dir2)

	files, _ := ioutil.ReadDir(targetPath)

	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			switch hierarchy {
			case Root:
				dir1 = fileName
			case First:
				dir2 = fileName
			case Other:
			}
			search(targetPath+"/"+fileName, dir1, dir2)
		} else {
			switch hierarchy {
			case Root:
				// TODO Rootのファイルの扱い
			case First:
				updateHierarchyMap(dir1, file.Size())
			case Other:
				updateHierarchyMap(dir1, file.Size())
				updateHierarchyMap(dir1+"/"+dir2, file.Size())
			}
		}
	}
}

// 階層を取得します
func getHierarchy(dir1 string, dir2 string) hierarchy {
	if dir1 == "" {
		return Root
	} else if dir2 == "" {
		return First
	} else {
		return Other
	}
}

// 階層容量情報を更新します
func updateHierarchyMap(key string, size int64) {
	if value, ok := hierarchyMap[key]; ok {
		value.Count++
		value.Size += size
		hierarchyMap[key] = value
	} else {
		hierarchyMap[key] = capacity{Size: size, Count: 1}
	}
}
