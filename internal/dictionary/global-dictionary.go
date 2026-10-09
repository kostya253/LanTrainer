package dictionary

import (
	"math/rand"
	"sync"

	dictionary_handler "github.com/kostya253/LanTrainer/internal/dictionary-handler"
)

var once_global sync.Once
var MyGloballDictionary *GeneralDictionary = nil

// Aka my dictionary
type GeneralDictionary struct {
	// a set of sentences
	sentences []string
}

func (p *GeneralDictionary) AddSentence(sentence string) {
	p.sentences = append(p.sentences, sentence)
}

func (p *GeneralDictionary) GetRandomSentence() string {
	return p.sentences[rand.Intn(len(p.sentences))]
}

func GetGlobalDictionaryInstance() DictionaryFunctor {
	once_global.Do(func() {
		// Read the my_dictionary.txt file and it's sentences
		my_sentences, err := dictionary_handler.ReadDictionary(GlobalDictionary)
		if err != nil {
			panic(err)
		}

		sentences := make([]string, len(my_sentences))
		copy(sentences, my_sentences)
		MyGloballDictionary = &GeneralDictionary{sentences: sentences}
	})

	return MyGloballDictionary
}
