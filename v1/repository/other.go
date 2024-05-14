package repository

import (
	"backend/v1/database"
	"backend/v1/model"
)

func CreateContact(contact *model.Contact) error {
	if err := database.DB.Create(contact).Error; err != nil {
		return err
	}

	return nil
}
