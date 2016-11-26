package main

import "io/ioutil"
import "fmt"
import "path/filepath"

// Hierarchy 階層を表す
type Hierarchy int

// 階層を表す
const (
	Root Hierarchy = iota
	First
	MoreSecond
)

func main() {
	rootPath := "C:/Workspace/golang/src/github.com/yuki-toida/hierarchy-search"
	search(rootPath, rootPath, "", "")
}

func search(rootPath string, targetPath string, dir1 string, dir2 string) {
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
			path := filepath.Join(targetPath, fileName)
			search(rootPath, path, dir1, dir2)
		} else {
			switch hierarchy {
			case Root:
				// TODO Rootのファイルの扱い不明
			case First:
				fmt.Println(dir1, file.Size())
			case MoreSecond:
				fmt.Println(dir2, file.Size())
			}
		}
	}
}

func getHierarchy(dir1 string, dir2 string) Hierarchy {
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
