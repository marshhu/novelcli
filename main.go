package main

import (
	"bytes"
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
	flag.Parse()

	if h {
		flag.Usage()
	}

	if g {
		//fmt.Fprintf(os.Stderr, "s:%s\n",s)
		//fmt.Fprintf(os.Stderr, "d:%s\n",d)
		if !strings.HasPrefix(s,"https://www.biquge.com.cn"){
			fmt.Fprint(os.Stderr, "Currently only supports generating novel text from Biquge(www.biquge.com.cn)\n")
			return
		}
		if len(d) <=0{
			fmt.Fprint(os.Stderr, "begin generate novel text...\n")
			return
		}
		fmt.Fprint(os.Stderr, "begin generate novel text...\n")
		novelService := service.NovelService{}
		novel, err := novelService.GetNovelByUrl(s)
		if err != nil {
			fmt.Fprint(os.Stderr, "crawler novel failed!")
		}
		var buffer bytes.Buffer
		for _, chapter := range novel.Chapters {
			buffer.WriteString(chapter.Name + "\n\n")
			buffer.WriteString(chapter.Content + "\n\n")
		}
		filePath := d + "/" + novel.Name + ".txt"
		err = utils.WriteFile(filePath, buffer.Bytes())
		if err != nil {
			fmt.Fprint(os.Stderr, "generate novel text failed!")
		}
		fmt.Fprintf(os.Stderr, "novel text path:%s", filePath)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `novlecli version: novlecli/1.0.0
      Usage: novlecli [-hg] [-s url] [-d path]
      Options:
     `)
	flag.PrintDefaults()
}
