package file

/*
递归遍历文件夹
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


/*
复制文件夹

base代表源目录,target代表目标目录
*/
func CopyDir(base string, target string) {
	os.MkdirAll(target, os.ModePerm)
	list := ReadDir(base)
	for _, f := range list {
		targetPath := strings.Replace(f, base, target, 1)
		TransFile(f, targetPath)
	}
}

// 上传目录到服务器
func UploadDir(base string, target string) bool {
	// list := ReadDir(base)
	return true
}

/*
	复制文件

source 源文件路径, target 目标文件路径

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

参数base是文件夹的路径

返回值是文件夹下所有的文件列表
*/
func ReadDir(base string) []string {
	list := make([]string, 0)
	ReadDirDeep(base, &list, nil)
	return list
}

/*
读取文件夹

参数base是文件夹的路径

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
