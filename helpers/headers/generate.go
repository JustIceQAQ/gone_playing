package headers

import "github.com/EDDYCJY/fake-useragent"

func Generate() string {
	return browser.Chrome()
}
