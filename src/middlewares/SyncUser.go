/*
 * SyncUser.go
 * Copyright (c) ti-bone 2023
 */

package middlewares

import (
	"errors"
	"feedbackBot/src/db"
	"feedbackBot/src/helpers"
	"feedbackBot/src/models"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gorm.io/gorm"
)

func SyncUser(b *gotgbot.Bot, ctx *ext.Context) error {
	// Work with DB in another goroutine, because handlerGroup is waiting until the function returns any value before procee
	go func() {
		var id = ctx.EffectiveMessage.From.Id
		var username = ctx.EffectiveMessage.From.Username
		var firstName = ctx.EffectiveMessage.From.FirstName
		var lastName = ctx.EffectiveMessage.From.LastName
		var languageCode = ctx.EffectiveMessage.From.LanguageCode

		var user models.User

		res := db.Connection.Where("user_id = ?", id).First(&user)

		user = models.User{
			UserID:       id,
			Username:     username,
			FirstName:    firstName,
			LastName:     lastName,
			LanguageCode: languageCode,
		}

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			resIns := db.Connection.Create(&user)
			if resIns.Error != nil {
				fmt.Printf("failed to insert user: %v", resIns.Error.Error())
			}

			err := helpers.LogUserAction("#NEW_USER\nNew user in the bot.", b, ctx)
			if err != nil {
				fmt.Printf("failed to send message to the logs: %v", err.Error())
			}
		} else {
			resUpd := db.Connection.Where("id = ?", user.ID).Updates(&user)
			if resUpd.Error != nil {
				fmt.Printf("failed to update user: %v", resUpd.Error.Error())
			}
		}
	}()

	return nil
}
