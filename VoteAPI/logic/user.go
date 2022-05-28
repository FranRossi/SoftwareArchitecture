package logic

import (
	"VoteAPI/data_access/repository"
	domain "VoteAPI/domain/user"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CheckVoter(idVoter string) (*domain.User, error) {
	user, err := repository.CheckVoterId(idVoter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, err
}

func RegisterUser(id string, username string, password string) (*domain.User, error) {
	user, err := createUser(id, username, password)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, err
}

func createUser(id string, username string, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	user := &domain.User{
		Id:             id,
		Username:       username,
		HashedPassword: string(hashedPassword),
	}
	return user, nil
}

func StoreUser(user *domain.User) (*domain.User, error) {
	err := repository.RegisterUser(user)
	if err != nil {
		return nil, fmt.Errorf("user cannot be created: %w", err)
	}
	return user, nil
}

func IsCorrectPassword(user *domain.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}
