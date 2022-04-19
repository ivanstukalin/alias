package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func getRandomWord(lang string) string {
	text := getText(lang)
	parsedText := parseText(text, lang)
	return parsedText
}

func getText(lang string) string {
	url := "https://random-word-api.herokuapp.com/word?lang=en"
	if lang == "RU" {
		url = "https://evilcoder.ru/random_word/"
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
	return string(body)
}

func parseText(text string, lang string) string {
	if lang == "EN" {
		return parseTextEn(text)
	}
	zp := regexp.MustCompile(`<tr><td>|</td><td>|</td></tr>`)
	words := zp.Split(text, -2)
	var parsedText []string
	for _, word := range words {
		if len(word) > 0 && !isWordSystem(word) {
			parsedText = append(parsedText, word)
		}
	}
	return parsedText[0]
}

func isWordSystem(word string) bool {
	systemWords := []string{
		"Сущ",
		"Прил",
		"Глаг",
	}
	for _, systemWord := range systemWords {
		if systemWord == word {
			return true
		}
	}
	return false
}

func parseTextEn(text string) string {
	parsedText := strings.Trim(text, "[\"")
	return strings.Trim(parsedText, "\"]")
}
