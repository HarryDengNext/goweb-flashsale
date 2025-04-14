package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"imooc_product/datamodels"
	"imooc_product/repositories"
	"log"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string)(user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User)(userId int64, err error)
}

func NewUserService(repository repositories.IUserRepository) IUserService  {
	return &UserService{repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string)(user *datamodels.User, isOk bool){
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		log.Println(err)
		return
	}
	isOk,_ = ValidatePassword(pwd, user.HashPassword)
	if !isOk {
		return &datamodels.User{}, false
	}
	return
}

func ValidatePassword(userPassword string, hashed string)(isOk bool, err error)  {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword));err != nil {
		return false, errors.New("密码比对错误")
	}
	return true,nil
}

func (u *UserService) AddUser(user *datamodels.User)(userId int64, err error){
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string)([]byte, error)  {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}
