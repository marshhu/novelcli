package service

import (
	"errors"
	"net/http"
	"net/url"
	"sort"
	"sync"

	"github.com/marshhu/novelcli/fetcher"
	"github.com/marshhu/novelcli/parser"
	"github.com/marshhu/novelcli/parser/biquge"
)

type NovelService struct {
}

func (s *NovelService) GetNovelByUrl(novelUrl string) (*Novel, error) {
	novel := Novel{}
	_, err := url.Parse(novelUrl)
	if err != nil {
		return nil, err
	}
	status, contents, err := fetcher.Fetcher(novelUrl, "", 60)
	if err != nil || status != http.StatusOK {
		return nil, errors.New("访问站点失败")
	}

	chapterParser := biquge.NewChapterListParser()
	chapterParseResult, err := chapterParser.Parse(novelUrl, contents)
	if err != nil {
		return nil, err
	}
	numOfWorkers := 30
	requestChan := make(chan string, numOfWorkers)
	resultChan := make(chan FetchResult, numOfWorkers)
	done := make(chan bool)
	var mutex sync.Mutex
	go createJob(chapterParseResult.Requests, requestChan)
	go handleResult(chapterParseResult.Requests, resultChan, &novel, &mutex, done)

	createWorkerPool(numOfWorkers, requestChan, resultChan)
	<-done
	data := chapterParseResult.Data.(parser.Novel)
	novel.FromModel(&data)
	sort.Sort(novel.Chapters)
	return &novel, nil
}

func createJob(requests map[string]parser.UrlParser, requestChan chan string) {
	for _, request := range requests {
		requestChan <- request.UrlInfo.Url
	}
	close(requestChan)
}

func createWorkerPool(numOfWorkers int, requestChan chan string, resultChan chan FetchResult) {
	var wg sync.WaitGroup
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go worker(requestChan, resultChan, &wg)
	}
	wg.Wait()
	close(resultChan)
}

func worker(requestChan <-chan string, resultChan chan FetchResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for request := range requestChan {
		status, contents, err := fetcher.Fetcher(request, "", 5)
		if err != nil || status != http.StatusOK {
			continue
		}
		fetchResult := FetchResult{Url: request, Content: contents}
		resultChan <- fetchResult
	}
}

func handleResult(requests map[string]parser.UrlParser, resultChan chan FetchResult, novel *Novel, mutex *sync.Mutex, done chan bool) {
	for fetchResult := range resultChan {
		if _, ok := requests[fetchResult.Url]; !ok {
			continue
		}
		urlParser := requests[fetchResult.Url]
		chapterParseResult, err := urlParser.Parse(fetchResult.Url, fetchResult.Content)
		if err != nil {
			continue
		}
		novelChapter := chapterParseResult.Data.(parser.NovelChapter)
		mutex.Lock()
		chapter := NovelChapter{Index: novelChapter.Index, Name: novelChapter.Name, Content: novelChapter.Content}
		novel.Chapters = append(novel.Chapters, chapter)
		mutex.Unlock()
	}
	done <- true
}
