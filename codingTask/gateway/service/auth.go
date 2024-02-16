package service

import (
	"context"
	"errors"
	"github.com/Nalivayko13/codingTask/gateway/model"
	"github.com/Nalivayko13/codingTask/gateway/utils"
	"net/http"
)

var UrlGenerateToken string

const (
	queryTokenGen = "login"
)

func (s *Service) AuthUser(ctx context.Context, user *model.User) (string, error) {
	//mock behavior of checking user in BD
	userDB := model.User{
		Password: "secretPass",
		Login:    "qwerty",
	}
	if userDB.Login != user.Login || userDB.Password != user.Password {
		return "", errors.New("no such user")
	}

	queryString := map[string]string{
		queryTokenGen: user.Login,
	}
	token, statsCode, err := utils.HttpGetCallWithParam(UrlGenerateToken, queryString, nil)
	if err != nil {
		return "", err
	}

	if statsCode != http.StatusOK {
		return "", errors.New("error in token generation")
	}

	return token, nil
}
