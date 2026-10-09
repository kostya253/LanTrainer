package botserver

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kostya253/LanTrainer/internal/dictionary"
	"github.com/kostya253/LanTrainer/internal/say"
	"github.com/kostya253/LanTrainer/internal/translator"
)

type BotServer struct {
	bot               *bot.Bot
	activeDictionary  string
	userLastSentences map[int64]string
	mu                sync.RWMutex
}

func NewBotServer() *BotServer {
	return &BotServer{
		activeDictionary:  dictionary.MyDictionary,
		userLastSentences: make(map[int64]string),
	}
}

func (b *BotServer) Start() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(b.handler),
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("BOT_TOKEN is not set")
	}

	bot, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.bot = bot
	bot.Start(ctx)
}

func (b *BotServer) handler(ctx context.Context, tbot *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	authorizedUsers := strings.Split(os.Getenv("AUTHORIZED_USER"), ",")
	isAuthorized := false
	for _, u := range authorizedUsers {
		if fmt.Sprintf("%d", userID) == strings.TrimSpace(u) {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Access denied.",
		})
		return
	}

	text := update.Message.Text
	if strings.HasPrefix(text, "/es:") {
		sentence := strings.TrimPrefix(text, "/es:")
		dictionary.GetDictionary(b.activeDictionary).AddSentence(sentence)
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Sentence added to " + b.activeDictionary,
		})
		return
	}

	if text == "/ts:my" {
		b.activeDictionary = dictionary.MyDictionary
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Switched to my dictionary",
		})
		return
	}

	if text == "/ts:g" {
		b.activeDictionary = dictionary.GlobalDictionary
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Switched to global dictionary",
		})
		return
	}

	b.mu.Lock()
	lastSentence := b.userLastSentences[userID]
	b.mu.Unlock()

	dict := dictionary.GetDictionary(b.activeDictionary)

	if lastSentence == "" {
		lastSentence = dict.GetRandomSentence()
		b.czReply(ctx, tbot, update, lastSentence)

		nextSentence := dict.GetRandomSentence()
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Translate to cz: %s", nextSentence),
		})
		b.mu.Lock()
		b.userLastSentences[userID] = nextSentence
		b.mu.Unlock()
	} else {
		translatedCZ, err := translator.TranslateUsingLLM(ctx, lastSentence)
		if err != nil {
			tbot.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Translation error: %v", err),
			})
			return
		}

		similarity := jaccardSimilarity(translatedCZ, text)
		similarityPercent := similarity * 100

		msg := "Not so much"
		if similarity > 0.5 {
			msg = "You are doing great!"
		}
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("%s score: %.2f", msg, similarityPercent),
		})

		b.czReply(ctx, tbot, update, lastSentence)

		newSentence := dict.GetRandomSentence()
		b.czReply(ctx, tbot, update, newSentence)

		nextSentence := dict.GetRandomSentence()
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Translate to cz: %s", nextSentence),
		})
		b.mu.Lock()
		b.userLastSentences[userID] = nextSentence
		b.mu.Unlock()
	}
}

func jaccardSimilarity(s1, s2 string) float64 {
	n := 2 // Using bigrams for partial matching
	getGrams := func(s string) map[string]struct{} {
		grams := make(map[string]struct{})
		s = strings.ToLower(s)
		runes := []rune(s)
		for i := 0; i+n <= len(runes); i++ {
			grams[string(runes[i:i+n])] = struct{}{}
		}
		return grams
	}

	set1 := getGrams(s1)
	set2 := getGrams(s2)

	if len(set1) == 0 && len(set2) == 0 {
		return 0
	}

	intersection := 0
	for gram := range set1 {
		if _, exists := set2[gram]; exists {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection
	return float64(intersection) / float64(union)
}

func (b *BotServer) czReply(ctx context.Context, tbot *bot.Bot, update *models.Update, text string) {
	czMsg, err := translator.TranslateUsingLLM(ctx, text)
	if err != nil {
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Translation error: %v", err),
		})
		return
	}

	if err := say.Say(czMsg); err != nil {
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Audio error: %v", err),
		})
		return
	}

	audioFile, err := os.Open("voice.ogg")
	if err != nil {
		tbot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("File error: %v", err),
		})
		return
	}
	defer audioFile.Close()

	tbot.SendAudio(ctx, &bot.SendAudioParams{
		ChatID:  update.Message.Chat.ID,
		Audio:   &models.InputFileUpload{Filename: "voice.ogg", Data: audioFile},
		Caption: fmt.Sprintf("en: %s\ncz: %s", text, czMsg),
	})

	os.Remove("voice.ogg")
	os.Remove("voice.wav")
}
