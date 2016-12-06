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
//const root = "C:/Workspace/Es/svn/"
const root = "\\\\gfs\\Shares\\00_全社共有\\"

// info 階層の容量を表す
type info struct {
	Size  int64
	Count int
}

// 階層容量マップ
var infoMap = map[string]info{}

// エントリポイント
func main() {
	start := time.Now()

	search()
	output()

	end := time.Now()
	fmt.Printf("%f(s)", (end.Sub(start)).Seconds())
}

// 階層探索します
func search() {
	fi, _ := ioutil.ReadDir(root)
	for _, f := range fi {
		if f.IsDir() {
			name := f.Name()
			infoMap[name] = info{Size: 0, Count: 0}
			search1(name)
		}
	}
}

// 第一階層
func search1(d1 string) {
	fi, _ := ioutil.ReadDir(root + "/" + d1)
	for _, f := range fi {
		if f.IsDir() {
			name := f.Name()
			infoMap[d1+"/"+name] = info{Size: 0, Count: 0}
			search2(d1, name)
		} else {
			update(d1, f.Size())
		}
	}
}

// 第二階層
func search2(d1 string, d2 string) {
	path := root + "/" + d1 + "/" + d2
	fi, _ := ioutil.ReadDir(path)
	for _, f := range fi {
		if f.IsDir() {
			searchN(d1, d2, path+"/"+f.Name())
		} else {
			size := f.Size()
			update(d1, size)
			update(d1+"/"+d2, size)
		}
	}
}

// その他階層
func searchN(d1 string, d2 string, path string) {
	fi, _ := ioutil.ReadDir(path)
	for _, f := range fi {
		if f.IsDir() {
			searchN(d1, d2, path+"/"+f.Name())
		} else {
			size := f.Size()
			update(d1, size)
			update(d1+"/"+d2, size)
		}
	}
}

func update(key string, size int64) {
	v, _ := infoMap[key]
	v.Count++
	v.Size += size
	infoMap[key] = v
}

// csvを出力します
func output() {
	dir, _ := os.Getwd()
	f, _ := os.Create(dir + "/hierarchy-search.csv")
	defer f.Close()
	//writer := bufio.NewWriterSize(f, 20000)
	w := bufio.NewWriterSize(transform.NewWriter(f, japanese.ShiftJIS.NewEncoder()), 20000)
	for k, v := range infoMap {
		w.WriteString(root + k + "," + strconv.FormatInt(v.Size, 10) + "," + strconv.Itoa(v.Count) + "\n")
	}
	w.Flush()
}

// func test() {
// 	currentDir, _ := os.Getwd()
// 	txt, _ := os.Open(currentDir + "/hoge.txt")
// 	defer txt.Close()
// 	scanner := bufio.NewScanner(txt)

// 	file, _ := os.Create(currentDir + "/hoge.csv")
// 	defer file.Close()
// 	writer := bufio.NewWriterSize(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()), 20000)

// 	for scanner.Scan() {
// 		content := scanner.Text()
// 		bytes, _ := writer.WriteString(content + "\n")
// 		fmt.Println(bytes)
// 	}
// 	writer.Flush()
// }
