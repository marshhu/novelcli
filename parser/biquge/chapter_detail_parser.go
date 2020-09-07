package biquge

import (
	"github.com/antchfx/htmlquery"
	"github.com/marshhu/novelcli/parser"
	"net/url"
	"regexp"
	"strings"
)

type ChapterDetailParser struct {
	novel *parser.Novel
}

func NewChapterDetailParser(_novel *parser.Novel) *ChapterDetailParser {
	return &ChapterDetailParser{novel: _novel}
}
func (p *ChapterDetailParser) Parse(crawlerUrl string, contents []byte) (*parser.ParseResult, error) {
	_, err := url.Parse(crawlerUrl)
	if err != nil {
		return nil, err
	}
	result := &parser.ParseResult{}
	result.Requests = make(map[string]parser.UrlParser)

	root, _ := htmlquery.Parse(strings.NewReader(string(contents)))
	findNode := htmlquery.FindOne(root, "//div[@id='content']")
	var chapterContent string
	if findNode != nil {
		chapterContent = htmlquery.InnerText(findNode)
		//替换空格&nbsp
		nbspRg := regexp.MustCompile(`&nbsp;`)
		chapterContent = nbspRg.ReplaceAllString(chapterContent, "")
		//替换换行<br>为\n
		brRg := regexp.MustCompile(`<br>`)
		chapterContent = brRg.ReplaceAllString(chapterContent, "\n")
		//fmt.Println(chapterContent)
		chapter := parser.NovelChapter{}
		chapter.Content = chapterContent
		if p.novel != nil {
			if _, ok := p.novel.Chapters[crawlerUrl]; ok {
				p.novel.Chapters[crawlerUrl].Content = chapterContent
				chapter.Index = p.novel.Chapters[crawlerUrl].Index
				chapter.Name = p.novel.Chapters[crawlerUrl].Name
			}
		}
		result.Data = chapter
	}

	return result, nil
}

func (p *ChapterDetailParser) Serialize() (name string, args interface{}) {
	return "biquge_chapter_detail_parser", nil
}
