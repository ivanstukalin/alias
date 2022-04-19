package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

var lang = ""
var skipButton = "SKIP"
var doneButton = "DONE"
var startButtonText = "Start"
var startText = "Click \"Start\" when you're ready to start explaining."
var endTimeText = "Время вышло, последнее слово!"
var changeLanguage = "Change language"
var answersCount = 0
var resultText = "Round result:"
var isGameInProcess = false

var langKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("RU"),
		tgbotapi.NewKeyboardButton("EN"),
	),
)

var gameButtons = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(skipButton),
		tgbotapi.NewKeyboardButton(doneButton),
	),
)

var startButton = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(startButtonText),
		tgbotapi.NewKeyboardButton(changeLanguage),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI("TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			text := ""

			switch update.Message.Text {
			case "/start":
				text = "Hello, please select a language"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ReplyMarkup = langKeyboard
				bot.Send(msg)
			case "RU", "EN":
				lang = update.Message.Text
				setSettingsByLang(lang)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, startText)
				msg.ReplyMarkup = startButton
				bot.Send(msg)
			case "Change language", "Изменить язык":
				sendLangCheck(bot, update.Message.Chat.ID)
			case "Start", "Начать":
				if len(lang) == 0 {
					sendLangCheck(bot, update.Message.Chat.ID)
					break
				}
				answersCount = 0
				sendRandomWord(bot, update.Message.Chat.ID)
				isGameInProcess = true
				time.AfterFunc(time.Second*60, func() {
					finishGame(bot, update.Message.Chat.ID)
				})
			case "DONE", "Угадал":
				answersCount++
				if !isGameInProcess {
					sendLastRoundResults(bot, update.Message.Chat.ID)
					sendNewGameInvite(bot, update.Message.Chat.ID)
					break
				}
				sendRandomWord(bot, update.Message.Chat.ID)
			default:
				if !isGameInProcess {
					sendLastRoundResults(bot, update.Message.Chat.ID)
					sendNewGameInvite(bot, update.Message.Chat.ID)
					break
				}
				sendRandomWord(bot, update.Message.Chat.ID)
			}
		}
	}
}

func setSettingsByLang(lang string) {
	if lang == "EN" {
		return
	}

	startText = "Нажмите \"Начать\" когда будете готовы начать объяснять"
	startButtonText = "Начать"
	skipButton = "Пропустить"
	doneButton = "Угадал"
	resultText = "Результат раунда:"
	changeLanguage = "Изменить язык"

	gameButtons = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(skipButton),
			tgbotapi.NewKeyboardButton(doneButton),
		),
	)

	startButton = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(startButtonText),
		),
	)

	startButton = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(startButtonText),
			tgbotapi.NewKeyboardButton(changeLanguage),
		),
	)
}

func sendRandomWord(bot *tgbotapi.BotAPI, chatId int64) {
	text := getRandomWord(lang)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = gameButtons
	bot.Send(msg)
}

func finishGame(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, endTimeText)
	bot.Send(msg)
	isGameInProcess = false
}

func sendLastRoundResults(bot *tgbotapi.BotAPI, chatId int64) {
	count := fmt.Sprintf("%d", answersCount)
	msg := tgbotapi.NewMessage(chatId, resultText+count)
	bot.Send(msg)
}

func sendNewGameInvite(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, startText)
	msg.ReplyMarkup = startButton
	bot.Send(msg)
}

func sendLangCheck(bot *tgbotapi.BotAPI, chatId int64) {
	text := "Please select a language (Пожалуйста, выберите язык)"
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = langKeyboard
	bot.Send(msg)
}
