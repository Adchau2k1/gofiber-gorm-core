package repository

import (
	"backend/v1/database"
	"backend/v1/model"
	"errors"

	"gorm.io/gorm"
)

func UserExists(id, username string) (bool, error) {
	db := database.DB
	var count int64
	var err error

	if id != "" {
		err = db.Model(&model.UserExport{}).Where("id = ?", id).Count(&count).Error
	} else {
		err = db.Model(&model.UserExport{}).Where("username = ?", username).Count(&count).Error
	}

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
func EmptyUsers() (bool, error) {
	db := database.DB
	result := db.Take(&model.UserExport{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return true, nil
	}

	return false, result.Error
}

func GetUser(id string, username string, page int, limit int, isAll bool) ([]model.UserExport, int, error) {
	db := database.DB
	users := []model.UserExport{}
	var err error

	if isAll {
		err = db.Order("created_at desc").Find(&users).Error
	} else {
		if page <= 0 {
			page = 1
		}
		if limit < 0 {
			limit = 0
		}

		offset := (page - 1) * limit
		if id != "" && username != "" {
			err = db.Limit(limit).Offset(offset).Where("id = ? AND username = ?", id, username).Order("created_at desc").Find(&users).Error
		} else if id != "" || username != "" {
			err = db.Limit(limit).Offset(offset).Where("id = ? OR username = ?", id, username).Order("created_at desc").Find(&users).Error
		} else {
			err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&users).Error
		}
	}

	if err != nil {
		return nil, 0, err
	}

	return users, len(users), nil
}
func GetTotalUser(id string, username string) (int64, error) {
	db := database.DB
	users := model.User{}
	var count int64
	var err error

	if id != "" && username != "" {
		sql := "id = ? AND username = ?"
		err = db.Model(&users).Where(sql, id, username).Count(&count).Error
	} else if id != "" || username != "" {
		sql := "id = ? OR username = ?"
		err = db.Model(&users).Where(sql, id, username).Count(&count).Error
	} else {
		err = db.Model(&users).Count(&count).Error
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CreateUser(user *model.User) error {
	if err := database.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func UpdateUser(id string, user *model.User) error {
	if err := database.DB.Model(&model.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return err
	}

	return nil
}

func DeleteUser(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return err
	}

	return nil
}
func ResetUser(listIds []string) error {
	if err := database.DB.Unscoped().Model(&model.User{}).Where("id IN (?)", listIds).Update("deleted_at", nil).Error; err != nil {
		return err
	}

	return nil
}
