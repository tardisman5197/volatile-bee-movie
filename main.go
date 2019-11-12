package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	inputFilePathPtr := flag.String("f", "", "A String containing the script.")

	flag.Parse()

	if len(*inputFilePathPtr) == 0 {
		fmt.Println("No script given")
		os.Exit(1)
	}

	script := readScript(*inputFilePathPtr)
	words := splitScript(script)
	sections := createSections(words)

	sendSections(sections)
}

func readScript(filepath string) string {
	scriptByt, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	return string(scriptByt)
}

func splitScript(script string) (words []string) {
	words = strings.Split(script, " ")
	return words
}

func createSections(words []string) (sections []*Section) {

	currentSectionID := 1

	sections = append(sections, newSection(currentSectionID))

	// Because section id starts at 1 when accessing sections - 1
	for i := 0; i < len(words); i++ {
		if sections[currentSectionID-1].checkAdd(words[i]) {
			sections[currentSectionID-1].addWord(words[i])
		} else {
			currentSectionID++
			sections = append(sections, newSection(currentSectionID))
			sections[currentSectionID-1].addWord(words[i])
		}
	}

	// last section should not point to another section
	sections[len(sections)-1].footer = ""

	return sections
}

func sendSections(sections []*Section) {
	for i := 0; i < len(sections); i++ {
		sections[i].join()

		sendData(sections[i].key, sections[i].body)
		time.Sleep(1 * time.Second)
		updateProgressbar(i, len(sections))
	}
}

func sendData(key, value string) {
	// fmt.Printf("key: %v, Value: %v\n", key, value)
	requestURL := fmt.Sprintf("%v?key=%v&val=%v", apiURL, key, value)
	// fmt.Printf("%v\n", requestURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatalf("It doesn't work: %v\n", err)
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Error: %v - %v\n", key, resp.Status)
		fmt.Printf("value: %v\n", value)
	}
}

func updateProgressbar(current, total int) {
	fmt.Printf("Progress %v / %v\r\n", current, total)
}
