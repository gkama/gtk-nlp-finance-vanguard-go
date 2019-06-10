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

//Model stores the model used to categorize the content
type Model struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Details  string  `json:"details"`
	Children []Model `json:"children"`
}

func main() {
	r := gin.Default()

	r.GET("/ping", ping)
	r.GET("/nlp/finance/vanguard/model", getModelJSON)
	r.POST("/nlp/finance/vanguard/categorize", categorize)

	r.Run()
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "Healthy")
}

func categorize(c *gin.Context) {
	var req ContentRequest

	c.BindJSON(&req)

	c.JSON(http.StatusOK, req)
}

func getModel() Model {
	return Model{
		ID:      "984ce69d-de79-478b-9223-ff6349514e19",
		Name:    "Vanguard",
		Details: "",
		Children: []Model{
			Model{
				ID:       "5ec6957d-4de7-4199-9373-d4a7fb59d6e1",
				Name:     "Index Funds",
				Details:  "vbiix|vbinx|vbisx|vbltx|vbmfx|vdaix|vdvix|veiex|veurx|vexmx|vfinx|vfsvx|vftsx|vfwix|vgovx|vgtsx|vhdyx|viaix|vigrx|vihix|vimsx|visgx|visvx|vivax|vlacx|vmgix|vmvix|vpacx|vtebx|vtibx|vtipx|vtsax|vtsmx|vtws",
				Children: []Model{},
			},
		},
	}
}
func getModelJSON(c *gin.Context) {
	model := getModel()
	c.JSON(http.StatusOK, model)
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
