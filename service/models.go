package service

import "github.com/marshhu/novelcli/parser"

type Novel struct {
	Name          string            `json:"name"`          //书名
	Image         string            `json:"image"`         //封面图
	Author        string            `json:"author"`        //作者
	Intro         string            `json:"intro"`         //简介
	Status        string            `json:"status"`        //状态  连载  完结
	LatestChapter string            `json:"latestChapter"` //最新章节
	UpdateTime    string            `json:"updateTime"`    //最近更新时间
	Chapters      ChapterCollection //章节
}

type NovelChapter struct {
	Index   int    `json:"index"`   //序号
	Name    string `json:"Name"`    //章节名
	Content string `json:"content"` //内容
}

type ChapterCollection []NovelChapter

//Len()
func (chapters ChapterCollection) Len() int {
	return len(chapters)
}

//Less()
func (chapters ChapterCollection) Less(i, j int) bool {
	return chapters[i].Index < chapters[j].Index
}

//Swap()
func (chapters ChapterCollection) Swap(i, j int) {
	chapters[i], chapters[j] = chapters[j], chapters[i]
}

type FetchResult struct {
	Url     string
	Content []byte
}

func (novel *Novel) FromModel(model *parser.Novel) {
	novel.Author = model.Author
	novel.Name = model.Name
	novel.Status = model.Status
	novel.Intro = model.Intro
	novel.Image = model.Image
	novel.UpdateTime = model.UpdateTime
	novel.LatestChapter = model.LatestChapter
}
