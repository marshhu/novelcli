package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/marshhu/novelcli/service"
	"github.com/marshhu/novelcli/utils"
	"os"
	"strings"
)

var (
	h bool
	g bool
	s string
	d string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&g, "g", false, "generate text novel")

	flag.StringVar(&s, "s", "https://www.biquge.com.cn/book/31833/", "novel url")
	flag.StringVar(&d, "d", ".", "text save path")
	// 改变默认的 Usage
	flag.Usage = usage
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	}()

	flag.Parse()

	if h {
		flag.Usage()
	}

	if g {
		//fmt.Fprintf(os.Stderr, "s:%s\n",s)
		//fmt.Fprintf(os.Stderr, "d:%s\n",d)
		if !strings.HasPrefix(s, "https://www.biquge.com.cn") {
			fmt.Fprintln(os.Stderr, "Currently only supports generating novel text from Biquge(www.biquge.com.cn)")
			return
		}
		if len(d) <= 0 {
			fmt.Fprintln(os.Stderr, "begin generate novel text...")
			return
		}

		fmt.Fprintln(os.Stderr, "begin generate novel text...")
		err := generateNovelText(s,d)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Fprintln(os.Stderr, "generate novel success,file path:"+d)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "novelcli version: novelcli/1.0.0\n"+
		"Usage: novelcli [-hg] [-s novel url] [-d save path]\n"+
		"Options:\n")
	flag.PrintDefaults()
}

//生成小说text文本
func generateNovelText(url, savePath string) error {
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
	filePath := savePath + "/" + novel.Name + ".txt"
	err = utils.WriteFile(filePath, buffer.Bytes())
	if err != nil {
		return errors.New("generate novel text failed")
	}
	return nil
}
