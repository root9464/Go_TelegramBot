package commands

import (
	"log"
	"root/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) addTask(update tgbotapi.Update) {
	state := &hb.state

	if update.Message.Text == "/addTask" || update.Message.Text == "Добавить ДЗ" {
		state.IsSubjectInput = true
		state.IsTaskInput = false
		state.SubjectName = ""
		state.Task = ""
		state.IsAddTask = true

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название предмета:")
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("Ошибка: не удалось отправить сообщение. \n", err)
		}
	} else if state.IsSubjectInput && !state.IsTaskInput && state.IsAddTask {
		state.SubjectName = update.Message.Text
		state.IsSubjectInput = false
		state.IsTaskInput = true

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите задание:")
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("Ошибка: не удалось отправить сообщение. \n", err)
		}
	} else if state.IsTaskInput {
		state.Task = update.Message.Text
		state.IsTaskInput = false

		go func() {
			homework := models.Homework{
				SubjectName: state.SubjectName,
				Task:        state.Task,
			}
			err := hb.db.Create(&homework).Error
			if err != nil {
				log.Println("Ошибка: не удалось добавить задание. \n", err)
			}
		}()

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно")
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("Ошибка: не удалось отправить сообщение. \n", err)
		}
	}
}
