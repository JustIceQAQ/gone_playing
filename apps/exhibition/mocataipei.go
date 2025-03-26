package exhibition

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gone_playing/helpers/crawler"
	"gone_playing/helpers/format"
	"gone_playing/helpers/headers"
	"gone_playing/helpers/storage"
	"gone_playing/helpers/translation"
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MocaTaipei struct {
	format.Base
}

func NewMocaTaipei() *MocaTaipei {
	return &MocaTaipei{
		Base: format.Base{
			Fullname: "台北當代藝術館",
			Code:     "MocaTaipei",
			Url:      "https://www.mocataipei.org.tw/tw/ExhibitionAndEvent",
			Items:    []*format.Exhibition{},
		},
	}
}

func (MocaTaipei *MocaTaipei) getResponse() ([]byte, error) {
	client := crawler.NewHttpClient()
	getRequest, err := client.NewGet(MocaTaipei.Url)
	if err != nil {
		return nil, fmt.Errorf("新請求產生失敗:", err)
	}
	getRequest.Header.Add("user-agent", headers.Generate())
	resp, err := client.Do(getRequest)
	defer client.Close()
	if err != nil {
		return nil, fmt.Errorf("請求失敗:", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("讀取 Body 失敗:", err)
	}
	defer client.Close()
	return data, nil
}

func (MocaTaipei *MocaTaipei) getParse(context []byte) ([]*goquery.Selection, error) {
	doc, err := translation.ParseHTML(io.NopCloser(bytes.NewReader(context)))
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 錯誤: %w", err)
	}
	var items []*goquery.Selection
	doc.Find("div.listFrameBox > div.list").Each(func(i int, s *goquery.Selection) {
		items = append(items, s)
	})
	return items, nil
}

func (MocaTaipei *MocaTaipei) getItem(item *goquery.Selection) *format.Exhibition {
	root_url := "https://www.mocataipei.org.tw"
	title := item.Find("h3.imgTitle").First().Text()
	source, _ := item.Find("a.textFrame").First().Attr("href")
	image, _ := item.Find("figure.imgFrame > img.img").First().Attr("src")

	var datetime []string
	item.Find("div.dateBox > div.date").Each(func(i int, s *goquery.Selection) {
		year := strings.TrimSpace(s.Find("span.year").First().Text())
		month_day := strings.ReplaceAll(strings.TrimSpace(s.Find("p.day").First().Clone().Children().Remove().End().Text()), " ", "")
		raw_date := strings.ReplaceAll(year+"/"+month_day, "/", "-")
		datetime = append(datetime, raw_date)
	})

	return &format.Exhibition{
		Title:    title,
		DateTime: strings.Join(datetime, " ~ "),
		Image:    root_url + image,
		Source:   source,
	}
}

func (MocaTaipei *MocaTaipei) Run() {
	body, err := MocaTaipei.getResponse()
	if err != nil {
		log.Fatal("Response 錯誤:", err)
	}
	items, err := MocaTaipei.getParse(body)
	if err != nil {
		log.Fatal("解析 HTML 錯誤:", err)
	}
	var datas []*format.Exhibition
	for _, item := range items {
		datas = append(datas, MocaTaipei.getItem(item))
	}

	mocataipei := NewMocaTaipei()
	mocataipei.Items = datas

	jsonData, err := json.MarshalIndent(mocataipei, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	storage.SaveToJson(mocataipei.Code+".json", jsonData)
}
