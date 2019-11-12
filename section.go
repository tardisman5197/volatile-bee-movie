package main

import (
	"fmt"
	"net/url"
	"strings"
)

// Section ...
type Section struct {
	key       string
	words     []string
	body      string
	sectionID int
	footer    string
}

func newSection(sectID int) *Section {
	words := make([]string, 0)
	s := &Section{
		words:     words,
		sectionID: sectID,
		key:       fmt.Sprintf("%v%d", keyFmt, sectID),
		footer:    fmt.Sprintf("\nNext Section: %v%d", keyFmt, sectID+1),
	}
	return s
}

func (s *Section) checkAdd(toAdd string) bool {
	s.join()

	if len(toAdd)+len(s.body) > limit {
		return false
	}

	return true
}

func (s *Section) join() {
	s.body = strings.Join(s.words[:], " ")
	s.body += s.footer
	s.body = url.QueryEscape(s.body)
}

func (s *Section) addWord(word string) {
	s.words = append(s.words, word)
}
