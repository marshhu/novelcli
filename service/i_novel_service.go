package service

type INovelService interface {
	GetNovelByUrl(url string) (*Novel, error)
}
