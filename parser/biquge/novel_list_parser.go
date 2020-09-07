package biquge

import (
	"github.com/antchfx/htmlquery"
	"github.com/marshhu/novelcli/parser"
	"net/url"
	"strings"
)

type NovelListParser struct {
}

func NewNovelListParser() *NovelListParser {
	return &NovelListParser{}
}

func (p *NovelListParser) Parse(crawlerUrl string, contents []byte) (*parser.ParseResult, error) {
	u, err := url.Parse(crawlerUrl)
	if err != nil {
		return nil, err
	}
	result := &parser.ParseResult{}
	result.Requests = make(map[string]parser.UrlParser)

	root, _ := htmlquery.Parse(strings.NewReader(string(contents)))
	list := htmlquery.Find(root, "//div[@id='hotcontent']/div/div/dl/dt")
	for _, row := range list {
		linkNode := htmlquery.FindOne(row, "./a")
		if linkNode != nil {
			link := htmlquery.SelectAttr(linkNode, "href")
			if !strings.HasPrefix(link, "http") && !strings.HasPrefix(link, "https") {
				link = u.Scheme + "://" + u.Host + link
			}
			result.Requests[link] = parser.UrlParser{Parser: NewChapterListParser(), UrlInfo: parser.UrlInfo{Url: link, Text: htmlquery.InnerText(linkNode)}}
		}
	}

	list = htmlquery.Find(root, "//div[@id='newscontent']/div/ul/li")
	for _, row := range list {
		spanNodes := htmlquery.Find(row, "./span")
		if len(spanNodes) > 0 {
			linkNode := htmlquery.FindOne(spanNodes[0], "./a")
			if linkNode != nil {
				link := htmlquery.SelectAttr(linkNode, "href")
				if !strings.HasPrefix(link, "http") && !strings.HasPrefix(link, "https") {
					link = u.Scheme + "://" + u.Host + link
				}
				result.Requests[link] = parser.UrlParser{Parser: NewChapterListParser(), UrlInfo: parser.UrlInfo{Url: link, Text: htmlquery.InnerText(linkNode)}}
			}
		}
	}

	return result, nil
}

func (p *NovelListParser) Serialize() (name string, args interface{}) {
	return "biquge_novel_list_parser", nil
}
