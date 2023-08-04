package main

import (
	"flag"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	currencyAPIURL = "https://www.cbr-xml-daily.ru/daily_json.js"
	currencyPair   = "RUB/KGS"
)

type currencyResponse struct {
	Valute map[string]struct {
		Value float64 `json:"Value"`
	} `json:"Valute"`
}

func getCurrencyRate() (float64, error) {
	resp, err := http.Get(currencyAPIURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data currencyResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	if currency, ok := data.Valute[currencyPair]; ok {
		return currency.Value, nil
	}

	return 0, fmt.Errorf("currency pair not found")
}

func sendNotification(chatID int64, message string, botToken string) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Error initializing Telegram bot:", err)
		return
	}

	msg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
}

func monitorCurrency(chatID int64, botToken string) {
	prevRate, err := getCurrencyRate()
	if err != nil {
		fmt.Println("Error fetching initial currency rate:", err)
		return
	}

	sendNotification(chatID, fmt.Sprintf("Monitoring %s has started. Current rate: %.2f RUB/KGS.", currencyPair, prevRate), botToken)

	ticker := time.NewTicker(60 * time.Second) // Периодичность проверки курса (в секундах), здесь установлено 60 секунд
	defer ticker.Stop()

	for range ticker.C {
		currentRate, err := getCurrencyRate()
		if err != nil {
			fmt.Println("Error fetching current currency rate:", err)
			continue
		}

		if currentRate != prevRate {
			sendNotification(chatID, fmt.Sprintf("%s rate has changed. Current rate: %.2f RUB/KGS.", currencyPair, currentRate), botToken)
			prevRate = currentRate
		}
	}
}

func main() {
	var botToken string
	flag.StringVar(&botToken, "token", "", "Telegram API token")
	flag.Parse()

	if botToken == "" {
		fmt.Println("Telegram API token not provided. Please set the -token flag.")
		return
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Current token:", botToken)
		fmt.Println("Error initializing Telegram bot:", err)
		return
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				sendNotification(update.Message.Chat.ID, "Hello! I'm monitoring RUB/KGS rate and will notify you about its changes.", botToken)
				go monitorCurrency(update.Message.Chat.ID, botToken)
			case "stop":
				sendNotification(update.Message.Chat.ID, "Monitoring RUB/KGS has been stopped.", botToken)
			}
		}
	}
}
