package exhibition

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gone_playing/helpers/crawler"
	"gone_playing/helpers/format"
	"gone_playing/helpers/headers"
	"gone_playing/helpers/storage"
	"gone_playing/helpers/translation"
	"io"
	"log"
)

type CKSMH struct {
	format.Base
}

func NewCKSMH() *CKSMH {
	return &CKSMH{
		Base: format.Base{
			Fullname: "中正紀念堂",
			Code:     "CKSMH",
			Url:      "https://www.cksmh.gov.tw/News_Actives_photo.aspx?n=6067&sms=14954",
			Items:    []*format.Exhibition{},
		},
	}
}

func (CKSMH *CKSMH) getResponse() ([]byte, error) {
	client := crawler.NewHttpClient()
	getRequest, err := client.NewGet(CKSMH.Url)
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

func (CKSMH *CKSMH) getParse(context []byte) ([]*goquery.Selection, error) {
	doc, err := translation.ParseHTML(io.NopCloser(bytes.NewReader(context)))
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 錯誤: %w", err)
	}
	var items []*goquery.Selection
	doc.Find("a.div-activity").Each(func(i int, s *goquery.Selection) {
		items = append(items, s)
	})
	return items, nil
}

func (CKSMH *CKSMH) getItem(item *goquery.Selection) *format.Exhibition {

	title, _ := item.First().Attr("title")
	address := item.Find("p.activity-season").Text()
	datetime := item.Find("p.activity-time").Text()
	image, _ := item.Find("img").Attr("src")
	source, _ := item.First().Attr("href")

	return &format.Exhibition{
		Title:    title,
		Address:  address,
		DateTime: datetime,
		Image:    image,
		Source:   source,
	}
}

func (CKSMH *CKSMH) Run() {
	body, err := CKSMH.getResponse()
	if err != nil {
		log.Fatal("Response 錯誤:", err)
	}
	items, err := CKSMH.getParse(body)
	if err != nil {
		log.Fatal("解析 HTML 錯誤:", err)
	}
	var datas []*format.Exhibition
	for _, item := range items {
		datas = append(datas, CKSMH.getItem(item))
	}

	cksmh := NewCKSMH()
	cksmh.Items = datas

	jsonData, err := json.MarshalIndent(cksmh, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	storage.SaveToJson(cksmh.Code+".json", jsonData)
}
