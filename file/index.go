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
