package dictionary

const (
	MyDictionary     = "my_dictionary.txt"
	GlobalDictionary = "dictionary.txt"
)

type DictionaryFunctor interface {
	// add a sentence to the dictionary
	AddSentence(sentence string)

	// get a random sentence from the dictionary
	GetRandomSentence() string
}

func GetDictionary(dictionaryType string) DictionaryFunctor {
	switch dictionaryType {
	case MyDictionary:
		return GetMyDictionaryInstance()
	case GlobalDictionary:
		return GetGlobalDictionaryInstance()
	default:
		return GetMyDictionaryInstance()
	}
}
