package extract

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/royalbhati/cssremove/scanner"
)

//CSS extracts the css for that particular file and save that in that location where the original file is located
func CSS(filep, csspath string, wg *sync.WaitGroup) {
	content, err := ioutil.ReadFile(csspath)
	if err != nil {
		log.Fatal(err)
	}

	requiredCSS, err := RequiredCSS(filep)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)
	scss := strings.ReplaceAll(text, "\n", "")

	a := scanner.New(scss)
	nextToken := a.Next()
	reqs := []string{}
	mapV := make(map[string]bool)
	temp := []string{}
	for nextToken.Type.String() != "EOF" {
		isElem, isValid, tv := validSelector(nextToken, requiredCSS, mapV)
		if isValid {
			temp = append(temp, tv)
			mapV[nextToken.Value] = true
			nextT := a.Next()
			for nextT.Value != "{" {
				bew := a.LookAhead()
				temp = append(temp, nextT.Value)

				mapV[nextT.Value] = true
				if bew == "{" && checkFor(mapV, requiredCSS, isElem, temp) {
					temp = nil
				}
				nextT = a.Next()
			}

			if len(temp) > 0 && validCSS(mapV, requiredCSS) {
				for nextT.Value != "}" {
					temp = append(temp, nextT.Value)
					mapV[nextT.Value] = true
					nextT = a.Next()
				}
				temp = append(temp, nextT.Value)
				reqs = append(reqs, temp...)
				mapV = make(map[string]bool)
				temp = nil
			}
		}
		nextToken = a.Next()
	}
	cssData := strings.Join(reqs, "")
	if err := saveFile(filep, cssData, wg); err != nil {
		log.Fatal(err)

	}
}

func checkFor(mapV map[string]bool, css []string, isElem bool, selectors []string) bool {
	for _, v := range css {
		if isElem && (selectors[0] != "." && selectors[0] != "#") {
			if selectors[0] == v || selectors[1] == v {
				return false
			}
		}
		if mapV[v] && !contains(v) {
			if selectors[0] != v && selectors[1] != v {
				dotRem := removeDot(selectors)
				prev := dotRem[0]
				for i, s := range dotRem {
					if prev == "," && s == v {
						return false
					}
					prev = dotRem[i]
				}
				return true
			}
			return false
		}

	}

	return true
}

func validCSS(mapV map[string]bool, css []string) bool {
	for _, v := range css {
		if mapV[v] {
			return true
		}
	}

	return false
}

func saveFile(filep, data string, wg *sync.WaitGroup) error {
	defer wg.Done()
	dirPath := filepath.Dir(filep)
	basePath := filepath.Base(filep)
	name := dirPath + "/" + strings.Split(basePath, ".")[0] + ".module" + ".css"
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	l, err := f.WriteString(data)
	if err != nil {
		f.Close()
		return err
	}
	fmt.Println(l, "bytes written successfully in", name)
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func validSelector(nextToken *scanner.Token, css []string, mapV map[string]bool) (bool, bool, string) {

	if nextToken.Value == "." || nextToken.Value == "#" {
		return false, true, nextToken.Value
	}
	if validElementSelector(nextToken.Value, css, mapV) {
		return true, true, nextToken.Value
	}

	return false, false, ""
}

func validElementSelector(token string, css []string, mapV map[string]bool) bool {
	var validID = regexp.MustCompile(`\w+\[\w+(\$|\*|)=".\w+"]`)
	for _, v := range css {
		if v == token {
			return true
		}
	}
	return validID.MatchString(token)
}

func removeDot(sels []string) []string {
	val := []string{}
	for i := range sels {
		if (sels)[i] == "." || (sels)[i] == " " {
			(sels)[i] = ""
		} else {
			val = append(val, sels[i])
		}
	}
	return val
}

func getElemnts() []string {
	return []string{"a",
		"abbr",
		"acronym",
		"address",
		"applet",
		"area",
		"article",
		"aside",
		"audio",
		"b",
		"base",
		"basefont",
		"bdi",
		"bdo",
		"bgsound",
		"big",
		"blink",
		"blockquote",
		"body",
		"br",
		"button",
		"canvas",
		"caption",
		"center",
		"cite",
		"code",
		"col",
		"colgroup",
		"content",
		"data",
		"datalist",
		"dd",
		"decorator",
		"del",
		"details",
		"dfn",
		"dir",
		"div",
		"dl",
		"dt",
		"element",
		"em",
		"embed",
		"fieldset",
		"figcaption",
		"figure",
		"font",
		"footer",
		"form",
		"frame",
		"frameset",
		"h1",
		"h2",
		"h3",
		"h4",
		"h5",
		"h6",
		"head",
		"header",
		"hgroup",
		"hr",
		"html",
		"i",
		"iframe",
		"img",
		"input",
		"ins",
		"isindex",
		"kbd",
		"keygen",
		"label",
		"legend",
		"li",
		"link",
		"listing",
		"main",
		"map",
		"mark",
		"marquee",
		"menu",
		"menuitem",
		"meta",
		"meter",
		"nav",
		"nobr",
		"noframes",
		"noscript",
		"object",
		"ol",
		"optgroup",
		"option",
		"output",
		"p",
		"param",
		"plaintext",
		"pre",
		"progress",
		"q",
		"rp",
		"rt",
		"ruby",
		"s",
		"samp",
		"script",
		"section",
		"select",
		"shadow",
		"small",
		"source",
		"spacer",
		"span",
		"strike",
		"strong",
		"style",
		"sub",
		"summary",
		"sup",
		"table",
		"tbody",
		"td",
		"template",
		"textarea",
		"tfoot",
		"th",
		"thead",
		"time",
		"title",
		"tr",
		"track",
		"tt",
		"u",
		"ul",
		"var",
		"video",
		"wbr",
		"xmp"}
}

func contains(elem string) bool {
	elems := getElemnts()
	for _, v := range elems {
		if elem == v {
			return true
		}
	}
	return false
}
