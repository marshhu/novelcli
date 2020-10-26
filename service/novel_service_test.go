package service

import (
	"testing"
)

func TestGetNovelByUrl(t *testing.T) {
	url := "https://www.biquge.com.cn/book/32146/"
	novelService := NovelService{}
	novel, err := novelService.GetNovelByUrl(url)
	if err != nil {
		t.FailNow()
	}
	if novel.Name != "开天录" && novel.Author != "血红" {
		t.FailNow()
	}
}
