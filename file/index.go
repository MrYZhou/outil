package file

import (
	"os"
	"path"
	"strings"
	"sync"

	. "github.com/MrYZhou/outil/common"
)

/*
合并文件

fileList 文件列表

target 文件输出路径
*/
func CombineFile(fileList []string, target string) {
	chunkTotal := make([]byte, 0)
	for _, name := range fileList {
		chunk, _ := os.ReadFile(name)
		chunkTotal = append(chunkTotal, chunk...)
	}
	os.WriteFile(target, []byte(chunkTotal), os.ModePerm)
}

/*
把一个文件切片成多个文件

out 切片输出的目录

filePath 切片的文件路径

num 切片数量
*/
func SliceFile(out string, filePath string, num int) []string {

	os.MkdirAll(out, os.ModePerm)

	f, _ := os.Open(filePath)
	fileInfo, _ := f.Stat()
	defer f.Close()

	size := fileInfo.Size() / int64(num)
	duo := fileInfo.Size() - size*int64(num)
	fileList := make([]string, 0)

	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {

			var chunk []byte
			if i == num-1 {
				chunk = make([]byte, size+duo)
			} else {
				chunk = make([]byte, size)
			}
			// 从源文件读取chunk大小的数据
			f.Read(chunk)

			rand_str := Random(10)

			targetPath := path.Join(out, "chunk"+rand_str)
			fileList = append(fileList, targetPath)

			os.WriteFile(targetPath, []byte(chunk), os.ModePerm)
			wg.Done()
		}(i)
	}

	wg.Wait()
	return fileList
}

/*
复制文件夹

base 代表源目录
target 代表目标目录
*/
func CopyDir(base string, target string) {
	os.MkdirAll(target, os.ModePerm)
	list := ReadDir(base)
	for _, f := range list {
		targetPath := strings.Replace(f, base, target, 1)
		TransFile(f, targetPath)
	}
}

/*
	复制文件

source 源文件路径, 

target 目标文件路径

注意中间的目录会自动创建,无需关心
*/
func TransFile(source string, target string) error {
	bytestr, _ := os.ReadFile(source)
	os.MkdirAll(path.Dir(target), 0755)
	err := os.WriteFile(target, []byte(bytestr), os.ModePerm)
	return err
}

/*
读取文件夹

base 文件夹的路径

返回值是文件夹下所有的文件列表
*/
func ReadDir(base string) []string {
	list := make([]string, 0)
	ReadDirDeep(base, &list, nil)
	return list
}

/*
读取文件夹

base 文件夹的路径

第一个返回值是文件夹下所有的文件列表

第二个返回值是文件夹列表
*/
func ReadDirAll(base string) ([]string, []string) {
	list := make([]string, 0)
	dirList := make([]string, 0)
	ReadDirDeep(base, &list, &dirList)
	return list, dirList
}

/*
递归遍历文件夹
base 文件夹的路径
list 文件列表
dirList 文件夹列表
*/
func ReadDirDeep(base string, list *[]string, dirList *[]string) {
	fileList, _ := os.ReadDir(base)
	for _, f := range fileList {
		if f.IsDir() {
			// 存下生成的目录
			if dirList != nil {
				*dirList = append(*dirList, path.Join(base, f.Name()))
			}
			ReadDirDeep(path.Join(base, f.Name()), list, dirList)
		} else {
			//  存下生成的文件
			*list = append(*list, path.Join(base, f.Name()))
		}
	}
}
