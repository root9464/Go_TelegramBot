package commands

import (
	"log"
	"root/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) addTask(update tgbotapi.Update) {
	if update.Message.Text == "/addTask" {
		subjectName, task = "", ""
		isSubjectInput = true
		isTaskInput = false
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите предмет:")
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("ошибка: не удалось отправить сообщение1. \n", err)
		}
	} else if isSubjectInput {
		subjectName = update.Message.Text
		isSubjectInput = false
		isTaskInput = true
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите задание:")
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("ошибка: не удалось отправить сообщение2. \n", err)
		}
	} else if isTaskInput {
		task = update.Message.Text
		homework := &models.Homework{
			SubjectName: subjectName,
			Task:        task,
		}
		if err := hb.db.Create(&homework).Error; err != nil {
			log.Println("ошибка: не удалось создать задание. \n", err)
		}
		confirmMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Задание добавлено!")
		if _, err := hb.bot.Send(confirmMsg); err != nil {
			log.Println("ошибка: не удалось отправить сообщение3. \n", err)
		}

		// Сброс значений
		subjectName = ""
		task = ""
		isSubjectInput = false
		isTaskInput = false
	}
}
