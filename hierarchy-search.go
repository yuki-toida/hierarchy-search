package main

import "io/ioutil"
import "fmt"
import "path/filepath"

// ルートディレクトリパス
const rootPath = "//gfs/Shares/00_全社共有/01_セキュリティ対策ソフト"

// 階層を表す
const (
	Root hierarchy = iota
	First
	MoreSecond
)

// Hierarchy 階層を表す
type hierarchy int16

// HierarchyInfo 階層情報を表す
type hierarchyInfo struct {
	Size  int64
	Count int32
}

// 階層マップ
var hierarchyMap = map[string]hierarchyInfo{}

// エントリポイント
func main() {
	search(rootPath, "", "")

	var csv string
	for k, v := range hierarchyMap {
		csv += fmt.Sprintf(k+",%v,%v\n", v.Size, v.Count)
	}
	fmt.Printf(csv)
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
			case MoreSecond:
			}
			nextPath := filepath.Join(targetPath, fileName)
			search(nextPath, dir1, dir2)
		} else {
			switch hierarchy {
			case Root:
				// TODO Rootのファイルの扱い
			case First:
				updateHierarchyMap(dir1, file.Size())
				//fmt.Println(dir1, file.Size())
			case MoreSecond:
				updateHierarchyMap(dir1, file.Size())
				updateHierarchyMap(dir1+"\\"+dir2, file.Size())
				//fmt.Println(dir1+"/"+dir2, file.Size())
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
		return MoreSecond
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
