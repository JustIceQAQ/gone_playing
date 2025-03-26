package translation

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
)

func ParseHTML(contest io.ReadCloser) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(contest)
	if err != nil {
		fmt.Println("解析 HTML 失敗:", err)
		return nil, fmt.Errorf("解析 HTML 失敗: %w", err)
	}
	return doc, nil

}
