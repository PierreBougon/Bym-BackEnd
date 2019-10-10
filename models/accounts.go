package models

import (
	u "github.com/PierreBougon/Bym-BackEnd/utils"

	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

/*
JWT claims struct
*/
type Token struct {
	UserId       uint
	TokenVersion uint
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	Model
	Email        string     `json:"email"`
	Password     string     `json:"password"`
	TokenVersion uint       `json:"token_version"`
	Picture      string     `json:"picture"`
	Playlists    []Playlist `gorm:"ForeignKey:UserId"`
}

func (account *Account) ValidatePassword() (map[string]interface{}, bool) {
	if len(account.Password) < 6 {
		return u.Message(false, "Password with more than 6 characters is required"), false
	}
	return nil, true
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	resp, passed := account.ValidatePassword()
	if !passed {
		return resp, false
	}
	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.TokenVersion = 0

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID, TokenVersion: account.TokenVersion}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["token"] = tokenString
	return response
}

func (account *Account) UpdateAccount() map[string]interface{} {
	retAcc := &Account{}
	err := db.First(&retAcc, account.ID).Error
	if err != nil /*|| playlist.UserId != user*/ {
		return u.Message(false, "Invalid account")
	}
	if account.Picture != "" {
		retAcc.Picture = account.Picture
	}
	db.Save(&retAcc)
	return u.Message(true, "Account successfully updated")
}

func (account *Account) DeleteAccount(user uint) map[string]interface{} {
	retAccount := &Account{}
/*	if user != userToDelete {
		return u.Message(false, "This account does not belong to you")
	}*/
	err := db.Table("accounts").Where("id = ?", user).First(&retAccount).Error
	// should not be possible since user is fetch from auth token
	if err != nil {
		return u.Message(false, "This account does not exist")
	}
	db.Delete(&retAccount)
	return u.Message(true, "Account successfully deleted")
}

func Login(email, password string) map[string]interface{} {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{UserId: account.ID, TokenVersion: account.TokenVersion}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	resp := u.Message(true, "Logged In")
	resp["token"] = tokenString
	return resp
}

func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}

func UpdatePassword(user uint, password string) map[string]interface{} {
	account := &Account{}
	db.First(account, user)

	resp, passed := account.ValidatePassword()
	if !passed {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	//Create new JWT token for the newly registered account
	account.TokenVersion++
	tk := &Token{UserId: account.ID, TokenVersion: account.TokenVersion}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	db.Save(account)

	ret := u.Message(true, "Password successfully updated")
	ret["token"] = tokenString
	return ret
}
