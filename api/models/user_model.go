package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint32 `gorm:"primary_key; auto_increment" json:"id"`
	FullName string `gorm:"size: 255;not_null;unique" json:"fullname"`
	Email    string `gorm:"size:255;not_null;unique" json:"email"`
	Password string `gorm:"size:100;not_null;" json:"password"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (user *User) BeforeSave() error {
	hashedPassword, err := Hash(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return nil

}

func (user *User) Prepare() {
	user.ID = 0
	user.FullName = html.EscapeString(strings.TrimSpace(user.FullName))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if user.FullName == "" {
			return errors.New("Required Name")
		}

		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Email == "" {
			return errors.New("Required Email")
		}

		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil

	case "login":
		if user.Password == "" {
			return errors.New("Required Password")
		}

		if user.Email == "" {
			return errors.New("Required Email")
		}

		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if user.FullName == "" {
			return errors.New("Required Nickname")
		}
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	}

}

func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	if err := db.Debug().Create(&user).Error; err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var users []User
	if err := db.Debug().Model(&User{}).Limit(100).Find(&users).Error; err != nil {
		return &[]User{}, err
	}
	return &users, nil
}

func (user *User) FindUserById(db *gorm.DB, id uint32) (*User, error) {
	if err := db.Debug().Model(User{}).Where("id = ?", id).Take(&user).Error; err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) UpdateUser(db *gorm.DB, id uint32) (*User, error) {
	if err := user.BeforeSave(); err != nil {
		log.Fatal(err)
	}

	db.Debug().Model(&User{}).Where("id = ?", id).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  user.Password,
			"fullName":  user.FullName,
			"email":     user.Email,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	if err := db.Debug().Model(&User{}).Where("id = ?", id).Take(&user).Error; err != nil {
		return &User{}, nil
	}

	return user, nil
}

func (user *User) DeleteUser(db *gorm.DB, id uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", id).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
