package dictionary_handler

import (
	"os"
	"strings"
)

func ReadDictionary(dictionary string) ([]string, error) {
	var sentences []string

	// Read the dictionary file and it's sentences
	data, err := os.ReadFile(dictionary)
	if err != nil {
		return nil, err
	}

	for _, sentence := range strings.Split(string(data), "\n") {
		sanitized := strings.TrimSpace(sentence) // removes \r, spaces, tabs, etc.
		if sanitized != "" {                     // guard against empty lines
			sentences = append(sentences, sanitized)
		}
	}

	return sentences, nil
}

func WriteDictionary(dictionary string, sentences []string) error {
	return nil
}

func WriteSentence(dictionary string, sentence string) error {
	// Append string to a dictionary file
	file, err := os.OpenFile(dictionary, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(sentence + "\n")

	if err != nil {
		return err
	}

	return nil
}
