package utils

import (
	"bufio"
	"fmt"
	"os"
)

func WriteFile(filePath string, content []byte) error {
	outputFile, outputError := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("打开或创建文件失败\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	outputWriter.Write(content)
	outputWriter.Flush()
	return nil
}

//  应用所在根目录
func RootDir() string {
	wd, _ := os.Getwd()
	return wd
}
