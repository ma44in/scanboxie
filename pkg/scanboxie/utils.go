package scanboxie

import "strings"

// Split ...
// Thy to https://ferencfbin.medium.com/golang-extended-split-function-6a34676e7b3c
func Split(str string, delimiter rune) []string {
	var isQuote = false
	f := func(c rune) bool {
		if c == '"' {
			if isQuote {
				isQuote = false
			} else {
				isQuote = true
			}
		}
		if !isQuote {
			return c == delimiter
		}
		return false

	}

	result := strings.FieldsFunc(str, f)

	for k, v := range result {
		if strings.Contains(v, `"`) {
			result[k] = strings.Replace(v, `"`, "", -1)
		}
	}

	return result
}
