package main

import (
	"log"
	"os"
	"root/database"
	"root/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	method := database.MethodDB()
	db, err := method.Connect()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных. \n", err)
	}

	var subjectName, task string
	var step int

	for update := range updates {
		if update.Message == nil || update.Message.From.ID == bot.Self.ID {
			continue
		}

		if update.Message.Text == "/task" {
			// Reset variables for the new homework entry
			subjectName = ""
			task = ""
			step = 0
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the subject name:")
			_, err := bot.Send(msg)
			if err == nil {
				step = 1
			}
		} else if step == 1 {
			subjectName = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the task:")
			_, err := bot.Send(msg)
			if err == nil {
				step = 2
			}
		} else if step == 2 {
			task = update.Message.Text
			homework := models.Homework{
				SubjectName: subjectName,
				Task:        task,
			}
			err := db.Create(&homework).Error
			if err == nil {
				confirmMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Homework saved successfully!")
				bot.Send(confirmMsg)
			}
			// Reset variables for the next homework entry
			subjectName = ""
			task = ""
			step = 0
		}
	}
}
