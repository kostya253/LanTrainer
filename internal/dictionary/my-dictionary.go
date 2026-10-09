package dictionary

import (
	"fmt"
	"log/slog"
	"math/rand"
	"sync"

	dictionary_handler "github.com/kostya253/LanTrainer/internal/dictionary-handler"
)

var once sync.Once
var MyPersonalDictionary *PersonalDictionary = nil

// Aka my dictionary
type PersonalDictionary struct {
	// a set of sentences
	sentences []string
}

func (p *PersonalDictionary) AddSentence(sentence string) {
	p.sentences = append(p.sentences, sentence)

	err := dictionary_handler.WriteSentence(MyDictionary, sentence)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to write sentence to global dictionary: %v", err))
	}
}

func (p *PersonalDictionary) GetRandomSentence() string {
	return p.sentences[rand.Intn(len(p.sentences))]
}

func GetMyDictionaryInstance() DictionaryFunctor {
	once.Do(func() {
		// Read the my_dictionary.txt file and it's sentences
		my_sentences, err := dictionary_handler.ReadDictionary(MyDictionary)
		if err != nil {
			panic(err)
		}

		sentences := make([]string, len(my_sentences))
		copy(sentences, my_sentences)
		MyPersonalDictionary = &PersonalDictionary{sentences: sentences}
	})

	return MyPersonalDictionary
}
