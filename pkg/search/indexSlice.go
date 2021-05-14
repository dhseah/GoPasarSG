package search

import (
	"strings"
	"sync"
	"unicode"

	"ProjectGoLive/pkg/models"

	sw "github.com/bbalet/stopwords"
	_ "github.com/go-sql-driver/mysql"
	snowballeng "github.com/kljensen/snowball/english"
)

// IndexSlice conains maps mapping words to a slice of integers
// containing ProductID(s) which has the word in their name,
// description or keywords fields
type IndexSlice []map[string][]int

// Add populates the IndexSlice by ranging through a list of
// Product and adding words from their name, decription, and
// keywords into their respective maps in the IndexSlice
func (idx IndexSlice) Add(products []*models.Product) {
	var wg sync.WaitGroup

	for _, product := range products {
		wg.Add(3)
		go func() {
			defer wg.Done()
			for _, token := range analyze(product.Name) {
				ids := idx[0][token]
				if ids != nil && ids[len(ids)-1] == product.ProductID {
					// Don't add same ID twice.
					continue
				}
				idx[0][token] = append(ids, product.ProductID)
			}
		}()
		go func() {
			defer wg.Done()
			for _, token := range analyze(product.Desc) {
				ids := idx[1][token]
				if ids != nil && ids[len(ids)-1] == product.ProductID {
					// Don't add same ID twice.
					continue
				}
				idx[1][token] = append(ids, product.ProductID)
			}
		}()
		go func() {
			defer wg.Done()
			for _, token := range analyze(product.Keyword) {
				ids := idx[2][token]
				if ids != nil && ids[len(ids)-1] == product.ProductID {
					// Don't add same ID twice.
					continue
				}
				idx[2][token] = append(ids, product.ProductID)
			}
		}()
		wg.Wait()
	}
}

// analyze prepares the text for search.
func analyze(text string) []string {
	text = stopword(text)            //Step 1
	tokens := tokenize(text)         //Step 2
	tokens = lowercaseFilter(tokens) //Step 3
	tokens = stemmerFilter(tokens)   //Step 4
	return tokens
}

// stopword removes stopwords to reduce false positives
// in the search result.
func stopword(text string) string {
	return sw.CleanString(text, "en", true)

}

// tokenize breaks down a phases and sentences
// into individual characters.
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// lowerCaseFilter converts every token to lower
// case so the search is not case-sensitive.
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// stemmerFilter stems every token so relevant results
// that do not match the search term exactly are include
// in the results.
func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}
