package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/marshhu/novelcli/service"
	"github.com/spf13/cobra"
)

var (
	url, output string
)

//初始化
func init() {
	txtCmd.Flags().StringVarP(&url, "url", "u", "",
		"novel web url (required)")
	txtCmd.MarkFlagRequired("url")
	txtCmd.Flags().StringVarP(&output, "output path", "o", "./",
		"output path")
}

var txtCmd = &cobra.Command{
	Use:   "txt",
	Short: "generate txt from novel web url",
	Long:  `generate txt from novel web url`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generateNovelText(); err != nil {
			return err
		}
		return nil
	},
}

//生成小说text文本
func generateNovelText() error {
	fmt.Println("generate novel begin...")
	novelService := service.NovelService{}
	novel, err := novelService.GetNovelByUrl(url)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("crawler novel failed")
	}
	var buffer bytes.Buffer
	for _, chapter := range novel.Chapters {
		buffer.WriteString(chapter.Name + "\n\n")
		buffer.WriteString(chapter.Content + "\n\n")
	}
	filePath := output + "/" + novel.Name + ".txt"
	if err := ioutil.WriteFile(filePath, buffer.Bytes(), os.ModePerm); err != nil {
		return errors.New("generate novel text failed")
	}
	return nil
}
