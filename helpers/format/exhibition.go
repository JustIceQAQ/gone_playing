package format

type Exhibition struct {
	Title    string `json:"title"`
	Address  string `json:"address"`
	DateTime string `json:"datetime"`
	Image    string `json:"image"`
	Source   string `json:"source"`
}

type Base struct {
	Fullname string        `json:"fullname"`
	Code     string        `json:"code"`
	Url      string        `json:"url"`
	Items    []*Exhibition `json:"item"`
}
