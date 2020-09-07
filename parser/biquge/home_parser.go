package biquge

import (
	"github.com/antchfx/htmlquery"
	"github.com/marshhu/novelcli/parser"
	"net/url"
	"strings"
)

type HomeParser struct {
}

func NewHomeParse() *HomeParser {
	return &HomeParser{}
}

func (p *HomeParser) Parse(crawlerUrl string, contents []byte) (*parser.ParseResult, error) {
	u, err := url.Parse(crawlerUrl)
	if err != nil {
		return nil, err
	}
	result := &parser.ParseResult{}
	result.Requests = make(map[string]parser.UrlParser)

	root, _ := htmlquery.Parse(strings.NewReader(string(contents)))
	list := htmlquery.Find(root, "//div[@class='nav']/ul/li")
	for _, row := range list {
		linkNode := htmlquery.FindOne(row, "./a")
		if linkNode != nil {
			link := htmlquery.SelectAttr(linkNode, "href")
			if link == "/" {
				continue
			}
			if !strings.HasPrefix(link, "http") && !strings.HasPrefix(link, "https") {
				link = u.Scheme + "://" + u.Host + link
			}
			result.Requests[link] = parser.UrlParser{Parser: NewNovelListParser(), UrlInfo: parser.UrlInfo{Url: link, Text: htmlquery.InnerText(linkNode)}}
		}
	}
	return result, nil
}

func (p *HomeParser) Serialize() (name string, args interface{}) {
	return "biquge_home_parser", nil
}
