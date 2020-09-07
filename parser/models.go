package parser

type Novel struct {
	Name          string                   `json:"name"`          //书名
	Image         string                   `json:"image"`         //封面图
	Author        string                   `json:"author"`        //作者
	Intro         string                   `json:"intro"`         //简介
	Status        string                   `json:"status"`        //状态  连载  完结
	LatestChapter string                   `json:"latestChapter"` //最新章节
	UpdateTime    string                   `json:"updateTime"`    //最近更新时间
	Chapters      map[string]*NovelChapter `json:"chapters"`      //章节
}

type NovelChapter struct {
	Index   int    `json:"index"`   //序号
	Name    string `json:"Name"`    //章节名
	Content string `json:"content"` //内容
}
