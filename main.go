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
type hierarchy int16

// HierarchyInfo 階層情報を表す
type hierarchyInfo struct {
	Size  int64
	Count int
}

// 階層マップ
var hierarchyMap = map[string]hierarchyInfo{}

// エントリポイント
func main() {
	start := time.Now()

	search(rootPath, "", "")
	output()

	end := time.Now()
	fmt.Printf("%f(s))", (end.Sub(start)).Seconds())
}

// csvを出力します
func output() {
	var csv string
	for k, v := range hierarchyMap {
		//csv += fmt.Sprintf(k+",%v,%v\n", v.Size, v.Count)
		csv += rootPath + k + "," + strconv.FormatInt(v.Size, 10) + "," + strconv.Itoa(v.Count) + "\n"
	}

	currentDir, _ := os.Getwd()
	err := ioutil.WriteFile(currentDir+"/hierarchy.csv", []byte(csv), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// currentDir, _ := os.Getwd()
	// file, _ := os.Create(currentDir + "/hierarchy.csv")
	// defer file.Close()
	// w := bufio.NewWriter(file)
	// w.WriteString(csv)
	// w.Flush()
}

// 階層を探索します
func search(targetPath string, dir1 string, dir2 string) {
	files, err := ioutil.ReadDir(targetPath)

	if err != nil {
		panic(err)
	}

	hierarchy := getHierarchy(dir1, dir2)

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
			nextPath := targetPath + "/" + fileName
			search(nextPath, dir1, dir2)
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
		// Root階層
		return Root
	} else if dir2 == "" {
		// 第一階層
		return First
	} else {
		// 第二階層以下
		return Other
	}
}

// 階層マップを更新します
func updateHierarchyMap(key string, size int64) {
	info, ok := hierarchyMap[key]
	if ok {
		info.Count++
		info.Size += size
		hierarchyMap[key] = info
	} else {
		hierarchyMap[key] = hierarchyInfo{Size: size, Count: 1}
	}
}
