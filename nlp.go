package main

import (
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

//ContentRequest is to receive the content to categorize
type ContentRequest struct {
	Content string `json:"content"`
}

func main() {
	r := gin.Default()

	r.GET("/ping", ping)
	r.POST("nlp/finance/vanguard/categorize", categorize)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "Healthy")
}

func categorize(c *gin.Context) {
	var req ContentRequest

	c.BindJSON(&req)

	c.JSON(http.StatusOK, req)
}

func tokenize(content string) []string {
	contentSplit := strings.FieldsFunc(content, split)
	stopWords := []string{"ourselves", "hers", "between", "yourself", "but", "again", "there", "about", "once", "during",
		"out", "very", "having", "with", "they", "own", "an", "be", "some", "for", "do", "its", "yours", "such",
		"into", "of", "most", "itself", "other", "off", "is", "s", "am", "or", "who", "as", "from", "him", "each",
		"the", "themselves", "until", "below", "are", "we", "these", "your", "his", "through", "don", "nor", "me",
		"were", "her", "more", "himself", "this", "down", "should", "our", "their", "while", "above", "both", "up",
		"to", "ours", "had", "she", "all", "no", "when", "at", "any", "before", "them", "same", "and", "been", "have",
		"in", "will", "on", "does", "yourselves", "then", "that", "because", "what", "over", "why", "so", "can", "did",
		"not", "now", "under", "he", "you", "herself", "has", "just", "where", "too", "only", "myself", "which", "those",
		"i", "after", "few", "whom", "t", "being", "if", "theirs", "my", "against", "a", "by", "doing", "it", "how",
		"further", "was", "here", "than"}

	return intersectSorted(contentSplit, stopWords).([]string)
}
func split(r rune) bool {
	return r == ' ' || r == ',' || r == ';' || r == '!' || r == '?' || r == '.'
}

func intersectSorted(a interface{}, b interface{}) interface{} {
	set := make([]interface{}, 0)
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		idx := sort.Search(bv.Len(), func(i int) bool {
			return bv.Index(i).Interface() == el
		})
		if idx < bv.Len() && bv.Index(idx).Interface() == el {
			set = append(set, el)
		}
	}

	return set
}
