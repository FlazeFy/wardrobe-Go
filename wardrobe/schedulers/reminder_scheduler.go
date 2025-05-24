package schedulers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SchedulerReminderUnansweredQuestion() {
	db := config.ConnectDatabase()

	// Get Admin Contact
	userContext := utils.NewUserContext(db)
	contact, err := userContext.GetAdminContact()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Get Unanswered Question
	questionContext := models.NewQuestionContext(db)
	questions, err := questionContext.GetUnansweredQuestion()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Send to Telegram
	datetime := time.Now()
	if len(contact) > 0 && len(questions) > 0 {
		filename := fmt.Sprintf("reminder_unanswered_question_%s.pdf", datetime)
		err = utils.GeneratePDFReminderUnansweredQuestion(questions, filename)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var list_question = ""
		for _, dt := range questions {
			list_question += fmt.Sprintf("- %s\nNotes: <i>ask at %s</i>\n\n", dt.Question, dt.CreatedAt)
		}

		for _, dt := range contact {
			if dt.TelegramUserId != nil && dt.TelegramIsValid {
				bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
				if err != nil {
					fmt.Println("Failed to connect to Telegram bot")
					return
				}

				telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
				if err != nil {
					fmt.Println("Invalid Telegram User Id")
					return
				}

				doc := tgbotapi.NewDocumentUpload(telegramID, filename)
				doc.ParseMode = "html"
				doc.Caption = fmt.Sprintf("[ADMIN] Hello %s, We're here to remind you. You have some unanswered question that needed to be answer. Here are the details:\n\n%s", dt.Username, list_question)

				_, err = bot.Send(doc)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}

		// Cleanup
		os.Remove(filename)
	}
}

func SchedulerReminderUnusedClothes() {
	db := config.ConnectDatabase()

	// Get Unused Clothes
	days := 60
	clothesContext := models.NewClothesContext(db)
	clothes, err := clothesContext.SchedulerGetUnusedClothes(days)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Send to Telegram
	if len(clothes) > 0 {
		var user_before = ""
		var list_clothes = ""

		for idx, dt := range clothes {
			if user_before == "" || user_before == dt.Username {

				var extra_desc = ""
				if dt.TotalUsed > 0 {
					extra_desc = fmt.Sprintf("Last used at %s", dt.LastUsed.Format("Y-m-d"))
				} else {
					extra_desc = "Never been used"
				}

				var extra_space = ""
				if idx < len(clothes)-1 {
					extra_space = "\n\n"
				}
				list_clothes += fmt.Sprintf("- %s (%s)\nNotes: <i>%s</i>%s", dt.ClothesName, dt.ClothesType, extra_desc, extra_space)
			}

			is_last := idx == len(clothes)-1
			is_different_user := !is_last && clothes[idx+1].Username != dt.Username

			if is_different_user || is_last {
				message := fmt.Sprintf("Hello %s, We're here to remind you. You have some clothes that has never been used since %d days after washed or being added to Wardrobe. Here are the details:\n\n%s\n\nUse and wash it again to keep your clothes at good quality and not smell musty", dt.Username, days, list_clothes)

				if dt.TelegramUserId != nil && dt.TelegramIsValid {
					bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
					if err != nil {
						fmt.Println("Failed to connect to Telegram bot")
						return
					}

					telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
					if err != nil {
						fmt.Println("Invalid Telegram User Id")
						return
					}

					doc := tgbotapi.NewMessage(telegramID, message)
					doc.ParseMode = "html"

					_, err = bot.Send(doc)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
				}

				list_clothes = ""
			}

			user_before = dt.Username
		}
	}
}

func SchedulerReminderUnironedClothes() {
	db := config.ConnectDatabase()

	// Get Unironed Clothes
	clothesContext := models.NewClothesContext(db)
	clothes, err := clothesContext.SchedulerGetUnironedClothes()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(clothes) > 0 {
		var user_before = ""
		var list_clothes = ""

		for idx, dt := range clothes {
			if user_before == "" || user_before == dt.Username {
				var extra_desc = " ("

				if dt.IsFavorite {
					extra_desc += "is your favorited"
				}
				if dt.IsScheduled {
					if dt.IsFavorite {
						extra_desc += ", "
					}
					extra_desc += "attached to schedule"
				}

				if dt.IsFavorite || dt.IsScheduled {
					extra_desc += ", "
				}
				if dt.HasWashed {
					extra_desc += "has"
				} else {
					extra_desc += "has'nt"
				}
				extra_desc += " been washed)"

				list_clothes += fmt.Sprintf("- %s%s\n<i>Notes: made from %s</i>\n", dt.ClothesName, extra_desc, dt.ClothesMadeFrom)
			}

			is_last := idx == len(clothes)-1
			is_different_user := !is_last && clothes[idx+1].Username != dt.Username

			if is_different_user || is_last {
				message := fmt.Sprintf("Hello %s, We're here to remind you. You have some clothes that has not been ironed yet. We only suggest the clothes that is made from cotton, linen, silk, or rayon. Here are the details:\n\n%s", dt.Username, list_clothes)

				// Send to Telegram
				if dt.TelegramUserId != nil && dt.TelegramIsValid {
					bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
					if err != nil {
						fmt.Println("Failed to connect to Telegram bot")
						return
					}

					telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
					if err != nil {
						fmt.Println("Invalid Telegram User Id")
						return
					}

					doc := tgbotapi.NewMessage(telegramID, message)
					doc.ParseMode = "html"

					_, err = bot.Send(doc)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
				}

				list_clothes = ""
			}

			user_before = dt.Username
		}
	}
}

func SchedulerReminderWashUsedClothes() {
	db := config.ConnectDatabase()
	days := 7

	// Get Unwashed Clothes
	clothesUsedContext := models.NewClothesUsedContext(db)
	clothes, err := clothesUsedContext.SchedulerGetUsedClothesReadyToWash(days)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(clothes) > 0 {
		var user_before = ""
		var list_clothes = ""

		for idx, dt := range clothes {
			if user_before == "" || user_before == dt.Username {
				var extra_desc = ""

				if dt.IsScheduled {
					extra_desc += "is on scheduled!"
				}
				if dt.IsFaded {
					if dt.IsScheduled {
						extra_desc += ", "
					}
					extra_desc += "is faded!"
				}

				if dt.IsScheduled || dt.IsFaded {
					extra_desc = fmt.Sprintf(", %s", extra_desc)
				}

				list_clothes += fmt.Sprintf("- <b>%s</b> (%s - %s)\n<i>Used Context: %s\nNotes: Last used at %s%s</i>\n\n", strings.Title(dt.ClothesName), strings.Title(dt.ClothesType), strings.Title(dt.ClothesMadeFrom), dt.UsedContext, dt.CreatedAt.Format("2006-01-02 15:04"), extra_desc)
			}

			is_last := idx == len(clothes)-1
			is_different_user := !is_last && clothes[idx+1].Username != dt.Username

			if is_different_user || is_last {
				message := fmt.Sprintf("Hello %s, We've noticed that some of your clothes are not washed after being used after %d days from now. Don't forget to wash your used clothes, here's the detail:\n\n%s", dt.Username, days, list_clothes)

				// Send to Telegram
				if dt.TelegramUserId != nil && dt.TelegramIsValid {
					bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
					if err != nil {
						fmt.Println("Failed to connect to Telegram bot")
						return
					}

					telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
					if err != nil {
						fmt.Println("Invalid Telegram User Id")
						return
					}

					doc := tgbotapi.NewMessage(telegramID, message)
					doc.ParseMode = "html"

					_, err = bot.Send(doc)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
				}

				list_clothes = ""
			}

			user_before = dt.Username
		}
	}
}
