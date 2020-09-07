package fetcher

import (
	"bufio"
	"errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func Fetcher(url string, cookie string, timeout int) (httpStatus int, body []byte, err error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return 0, nil, errors.New("url格式错误")
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	userAgent := "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"
	req.Header.Add("User-Agent", userAgent)

	if len(cookie) > 0 {
		req.Header.Add("cookie", cookie)
	}

	resp, error := client.Do(req)
	if error != nil || resp == nil {
		return 0, nil, errors.New("获取失败")
	}

	bodyReader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	e, _ := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	resBody, _ := ioutil.ReadAll(utf8Reader)
	return resp.StatusCode, resBody, nil
}

func determineEncoding(r *bufio.Reader) (encoding.Encoding, string) {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("fetcher error: %v", err)
		return unicode.UTF8, "utf8"
	}
	e, name, _ := charset.DetermineEncoding(bytes, "")
	if name == "windows-1252" {
		return simplifiedchinese.GBK, "gbk"
	}
	return e, name
}

func newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	userAgent := "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"
	req.Header.Add("User-Agent", userAgent)

	return req, err
}
