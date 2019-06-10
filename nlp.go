package main

import (
	"net/http"
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

//Category is used to return the categorization result
type Category struct {
	Name        string    `json:"name"`
	TotalWeight int       `json:"total_weight"`
	Matched     []Matched `json:"matched"`
}

//Matched is used to store the words/phrases matched in the categorization result
type Matched struct {
	Value  string `json:"value"`
	Weight int    `json:"weight"`
}

func main() {
	r := gin.Default()

	r.POST("/nlp/finance/vanguard/categorize", categorize)
	r.GET("/nlp/finance/vanguard/model", getModelJSON)
	r.GET("/nlp/finance/vanguard/ping", ping)

	r.Run()
}

func categorize(c *gin.Context) {
	var req ContentRequest
	var categories []Category

	c.BindJSON(&req)

	contentTokenized := tokenize(req.Content)
	model := getModel()
	modelStack := new(Stack)

	for _, mc := range model.Children {
		modelStack.push(mc)
	}

	for modelStack.len() > 0 {
		p := modelStack.pop().(Model)

		//Binary search on each token in content
		for _, tc := range contentTokenized {
			if contains(strings.Split(p.Details, "|"), tc) {
				categories = addCategory(categories, p, tc)
			}
		}

		//Add each child to the stack if it exists
		for _, c := range p.Children {
			modelStack.push(c)
		}
	}

	c.JSON(http.StatusOK, categories)
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
	c.JSON(http.StatusOK, getModel())
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

	//Remove stop words
	for i, cs := range contentSplit {
		if contains(stopWords, cs) {
			contentSplit = append(contentSplit[:i], contentSplit[i+1:]...)
		}
	}

	return contentSplit
}
func split(r rune) bool {
	return r == ' ' || r == ',' || r == ';' || r == '!' || r == '?' || r == '.'
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "Healthy")
}

//Stack to implement LIFO object and breadth first search
type Stack struct {
	top  *item
	size int
}
type item struct {
	value interface{} //value as interface type to hold any data type
	next  *item
}

func (stack *Stack) len() int {
	return stack.size
}
func (stack *Stack) push(value interface{}) {
	stack.top = &item{
		value: value,
		next:  stack.top,
	}
	stack.size++
}
func (stack *Stack) pop() (value interface{}) {
	if stack.len() > 0 {
		value = stack.top.value
		stack.top = stack.top.next
		stack.size--
		return
	}

	return nil
}

//Add category if it doesn't exist
//else, you update the matched weight
func addCategory(categories []Category, model Model, value string) []Category {
	//check if categories has the model.Name
	var cat Category
	var catIndex int
	for i, c := range categories {
		if c.Name == model.Name {
			cat = c
			catIndex = i
		}
	}

	//if cat is not empty then add a new one
	//else if cat is not empty and there is a matched, update the weight
	if cat.Name == "" {
		categories = append(categories, Category{
			Name:        model.Name,
			TotalWeight: 1,
			Matched: []Matched{
				Matched{
					Value:  value,
					Weight: 1,
				},
			},
		})
	} else {
		var matched Matched
		for i, m := range cat.Matched {
			if m.Value == value {
				matched = m

				categories[catIndex].TotalWeight++
				categories[catIndex].Matched[i].Weight++
			}
		}

		//if matched.Value is empty, then add it
		//else update the weight
		if matched.Value == "" {
			categories[catIndex].TotalWeight++
			categories[catIndex].Matched = append(categories[catIndex].Matched, Matched{
				Value:  value,
				Weight: 1,
			})
		}
	}

	return categories
}

//Helper functions
func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
