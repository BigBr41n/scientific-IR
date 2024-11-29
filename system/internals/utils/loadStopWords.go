package utils

import (
	"bufio"
	"os"
)

// LoadStopWords reads a .txt file and returns a map of stop words
func LoadStopWords() (map[string]struct{}, error) {
	stopWords := make(map[string]struct{})

	file, err := os.Open("./data/stop-words.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		stopWords[word] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stopWords, nil
}