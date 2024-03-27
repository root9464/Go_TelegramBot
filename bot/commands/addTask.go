package commands

import (
	"log"
	"root/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) addTask(update tgbotapi.Update) {
	if update.Message.Text == "/task" {
		subjectName, task = "", ""
		isSubjectInput = true
		isTaskInput = false
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите предмет:")
		_, err := hb.bot.Send(msg)
		if err == nil {
			log.Fatal("ошибка: не удалось отправить сообщение. \n", err)
		}
	} else if isSubjectInput {
		subjectName = update.Message.Text
		isSubjectInput = false
		isTaskInput = true
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите задание:")
		_, err := hb.bot.Send(msg)
		if err == nil {
			log.Fatal("ошибка: не удалось отправить сообщение. \n", err)
		}
	} else if isTaskInput {
		task = update.Message.Text
		homework := &models.Homework{
			SubjectName: subjectName,
			Task:        task,
		}
		err := hb.db.Create(&homework).Error
		if err != nil {
			log.Fatal("ошибка: не удалось создать задание. \n", err)
		}
		confirmMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Задание добавлено!")
		if _, err := hb.bot.Send(confirmMsg); err != nil {
			log.Fatal("ошибка: не удалось отправить сообщение. \n", err)
		}

		// Сброс значений
		subjectName = ""
		task = ""
		isSubjectInput = false
		isTaskInput = false
	}
}
