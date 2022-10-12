package utils

import (
	"bufio"
	"os"
)

func LoadFile(filename string) (lines []string) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close() //nolint
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	return
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
