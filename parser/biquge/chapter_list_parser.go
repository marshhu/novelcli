package biquge

import (
	"github.com/antchfx/htmlquery"
	"github.com/marshhu/novelcli/parser"
	"net/url"
	"regexp"
	"strings"
)

type ChapterListParser struct {
}

func NewChapterListParser() *ChapterListParser {
	return &ChapterListParser{}
}

func (p *ChapterListParser) Parse(crawlerUrl string, contents []byte) (*parser.ParseResult, error) {
	u, err := url.Parse(crawlerUrl)
	if err != nil {
		return nil, err
	}
	result := &parser.ParseResult{}
	result.Requests = make(map[string]parser.UrlParser)

	root, _ := htmlquery.Parse(strings.NewReader(string(contents)))
	//获取小说信息
	novel := parser.Novel{}
	novel.Chapters = make(map[string]*parser.NovelChapter)

	findNode := htmlquery.FindOne(root, "//div[@id='info']")
	if findNode != nil {
		h1Node := htmlquery.FindOne(findNode, "./h1")
		if h1Node != nil {
			novel.Name = htmlquery.InnerText(h1Node)
		}
		pNodes := htmlquery.Find(findNode, "./p")
		if len(pNodes) >= 4 {
			authorText := htmlquery.InnerText(pNodes[0])
			if len(authorText) > 0 {
				array := strings.Split(authorText, "：")
				if len(array) >= 2 {
					novel.Author = array[1]
				}
			}
			statusText := htmlquery.InnerText(pNodes[1])
			if len(statusText) > 0 {
				array := strings.Split(statusText, "：")
				if len(array) >= 2 {
					list := strings.Split(array[1], ",")
					if len(list) >= 1 {
						novel.Status = list[0]
					}
				}
			}
			novel.UpdateTime = htmlquery.InnerText(pNodes[2])
			novel.LatestChapter = htmlquery.InnerText(pNodes[3])
		}
	}
	findNode = htmlquery.FindOne(root, "//div[@id='intro']")
	if findNode != nil {
		introText := htmlquery.InnerText(findNode)
		rg := regexp.MustCompile(`\n`)
		introText = rg.ReplaceAllString(introText, "")
		rg = regexp.MustCompile(`\t`)
		introText = rg.ReplaceAllString(introText, "")
		rg = regexp.MustCompile(` `)
		introText = rg.ReplaceAllString(introText, "")
		novel.Intro = introText
	}
	findNode = htmlquery.FindOne(root, "//div[@id='fmimg']/img")
	if findNode != nil {
		novel.Image = htmlquery.SelectAttr(findNode, "src")
	}

	list := htmlquery.Find(root, "//div[@id='list']/dl/dd")
	for index, row := range list {
		linkNode := htmlquery.FindOne(row, "./a")
		if linkNode != nil {
			link := htmlquery.SelectAttr(linkNode, "href")
			if !strings.HasPrefix(link, "http") && !strings.HasPrefix(link, "https") {
				link = u.Scheme + "://" + u.Host + link
			}
			novelChapter := parser.NovelChapter{Index: index + 1, Name: htmlquery.InnerText(linkNode)}
			novel.Chapters[link] = &novelChapter
			result.Requests[link] = parser.UrlParser{Parser: NewChapterDetailParser(&novel), UrlInfo: parser.UrlInfo{Url: link, Text: htmlquery.InnerText(linkNode)}}
		}
	}
	//fmt.Println(chapters)
	result.Data = novel
	return result, nil
}

func (p *ChapterListParser) Serialize() (name string, args interface{}) {
	return "biquge_chapter_list_parser", nil
}
