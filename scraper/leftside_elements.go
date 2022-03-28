package scraper

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func getLeftsideElements(contentWrapper *colly.HTMLElement) LeftsideElements {
	leftsideElements := make([]*colly.HTMLElement, 0)
	contentWrapper.ForEach("div.leftside>div.spaceit_pad", func(i int, el *colly.HTMLElement) {
		leftsideElements = append(leftsideElements, el)
	})
	for i, j := 0, len(leftsideElements)-1; i < j; i, j = i+1, j-1 {
		leftsideElements[i], leftsideElements[j] = leftsideElements[j], leftsideElements[i]
	}
	return leftsideElements
}

type LeftsideElements []*colly.HTMLElement

func (lse LeftsideElements) Check() bool {
	return len(lse) >= 20
}

func (lse LeftsideElements) _getElement(el *colly.HTMLElement, extraText ...string) string {
	if tmp := el.ChildAttr("a[title]", "title"); tmp != "" {
		return tmp
	}
	extraText = append(extraText, "\n")
	return RemoveAndTrim(el.Text, extraText...)
}

func (lse LeftsideElements) getElement(element string, extraText ...string) string {
	// idx = len(lse) - idx - 1
	for _, el := range lse {
		if strings.Contains(el.Text, element+":") {
			extraText = append(extraText, element+":")
			return lse._getElement(el, extraText...)
		}
	}
	// if lse.Check() {
	// 	if tmp := lse[idx].ChildAttr("a[title]", "title"); tmp != "" {
	// 		return tmp
	// 	}
	// 	extraText = append(extraText, "\n")
	// 	return RemoveAndTrim(lse[idx].Text, extraText...)
	// }
	// log.Println("Error: invalid index ", idx)
	return ""
}

func (lse LeftsideElements) getElementList(element string, extraText ...string) []string {
	for _, el := range lse {
		if strings.Contains(el.Text, element+":") {
			if tmp := el.ChildAttrs("a[title]", "title"); len(tmp) != 0 {
				return tmp
			}
			extraText = append(extraText, element+":")
			return SplitAndTrim(lse._getElement(el, extraText...), ",")
		}
	}
	// // idx = len(lse) - idx - 1
	// if lse.Check() {
	// 	if tmp := lse[idx].ChildAttrs("a[title]", "title"); len(tmp) != 0 {
	// 		return tmp
	// 	}
	// 	return SplitAndTrim(lse.getElement(idx, extraText...), ",")
	// }
	// // log.Println("Error: invalid index ", idx)
	return []string{}
}

func (lse LeftsideElements) getElementInt(element string) int {
	for _, el := range lse {
		if strings.Contains(el.Text, element+":") {
			n, err := strconv.Atoi(RemoveAndTrim(el.Text, element+":", ",", "#", "\n"))
			if err != nil {
				log.Println("\033[31mError\033[0m ", err)
				return -1
			}
			return n
		}
	}
	return -1
}

func (lse LeftsideElements) Favorites() int {
	return lse.getElementInt("Favorites")
}

func (lse LeftsideElements) Members() int {
	return lse.getElementInt("Members")
}

func (lse LeftsideElements) Popularity() int {
	return lse.getElementInt("Popularity")
}

func (lse LeftsideElements) Ranked() int {
	tmp_l := strings.Split(lse.getElement("Ranked", "#"), " 2 ")
	if len(tmp_l) <= 0 {
		log.Println("\033[31mError\033[0m ", errors.New("invalid ranked text"))
		return 0
	}
	ranked, err := strconv.Atoi(strings.TrimSpace(tmp_l[0]))
	if err != nil {
		log.Println("\033[31mError\033[0m ", errors.New("invalid ranked text '"+strings.TrimSpace(tmp_l[0])+"'"))
		return 0
	}
	return IntWithoutLastChar(ranked)
}

func (lse LeftsideElements) Score() (float32, int) {
	scoreTxt := lse[4].ChildText("span[itemprop=ratingValue]")
	scoredByTxt := lse[4].ChildText("span[itemprop=ratingCount]")
	score, err := strconv.ParseFloat(scoreTxt, 32)
	if err != nil {
		log.Println("Score error: ", err)
	}
	scoredBy, err := strconv.Atoi(scoredByTxt)
	if err != nil {
		log.Println("Scored by error: ", err)
	}
	return float32(score), scoredBy
}

func (lse LeftsideElements) Rating() string {
	return lse.getElement("Rating")
}

func (lse LeftsideElements) Duration() string {
	return lse.getElement("Duration")
}

func (lse LeftsideElements) Demographic() string {
	return lse.getElement("Demographic")
}

func (lse LeftsideElements) Theme() string {
	return lse.getElement("Theme")
}

func (lse LeftsideElements) Genres() []string {
	return lse.getElementList("Genres")
}

func (lse LeftsideElements) Source() string {
	return lse.getElement("Source")
}

func (lse LeftsideElements) Studios() []string {
	return lse.getElementList("Studios")
}

func (lse LeftsideElements) Licensors() []string {
	return lse.getElementList("Licensors")
}

func (lse LeftsideElements) Producers() []string {
	return lse.getElementList("Producers")
}

func (lse LeftsideElements) Broadcast() string {
	return lse.getElement("Broadcast")
}

func (lse LeftsideElements) Premiered() string {
	return lse.getElement("Premiered")
}

func (lse LeftsideElements) Aired() string {
	return lse.getElement("Aired")
}

func (lse LeftsideElements) Status() string {
	return lse.getElement("Status")
}

func (lse LeftsideElements) Episodes() int {
	return lse.getElementInt("Episodes")
}

func (lse LeftsideElements) Type() string {
	return lse.getElement("Type")
}

func RemoveAndTrim(s string, old ...string) string {
	for _, v := range old {
		s = strings.ReplaceAll(s, v, "")
	}
	return strings.TrimSpace(s)
}

func SplitAndTrim(s, sep string) []string {
	tmp := make([]string, 0)
	for _, v := range strings.Split(s, sep) {
		tmp = append(tmp, strings.TrimSpace(v))
	}
	return tmp
}

func WithoutLast(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 0 {
		s = s[:len(s)-1]
	}
	return s
}

func IntWithoutLastChar(i int) int {
	s := strconv.Itoa(i)
	if len(s) <= 1 {
		return i
	}
	ni, err := strconv.Atoi(WithoutLast(s))
	if err != nil {
		return i
	}
	return ni
}
