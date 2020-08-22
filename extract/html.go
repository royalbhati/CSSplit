package extract

import (
	"bytes"
	"io"
	"io/ioutil"

	"golang.org/x/net/html"
)

//RequiredCSS extracts class names and elements from the specified files
func RequiredCSS(filep string) ([]string, error) {
	htmlbody, err := ioutil.ReadFile(filep)
	if err != nil {
		return nil, err
	}
	requiredCSS := []string{}
	htmlText := string(htmlbody)
	buff := bytes.NewBufferString(htmlText)
	tokenizer := html.NewTokenizer(buff)
	hElem := make(map[string]bool)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
		}
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			if !hElem[token.Data] {
				requiredCSS = append(requiredCSS, token.Data)
			} else {
				hElem[token.Data] = true
			}
			if token.Attr != nil {
				for _, v := range token.Attr {
					if v.Key == "class" || v.Key == "id" || v.Key == "classname" {
						requiredCSS = append(requiredCSS, v.Val)
					}
				}
			}
		}
	}
	return requiredCSS, nil
}
