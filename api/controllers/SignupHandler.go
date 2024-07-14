package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/Iyusuf40/goBackendUtils/mail"
	"github.com/Iyusuf40/goBackendUtils/models"
	"github.com/Iyusuf40/goBackendUtils/storage"
	"github.com/google/uuid"
)

type SignupHandler struct {
	temp_store  storage.TempStore
	users_store storage.Storage[models.User]
}

const DEFAULT_SIGNUP_TIMEOUT = 86400.0

func (sighnup_h *SignupHandler) HandleSignup(user models.User) string {

	userJson, err := json.Marshal(user)

	if err != nil {
		return ""
	}

	signupId := uuid.NewString()

	userEmail := user.Email

	go func() {
		sighnup_h.sendEmailConfirmationMsg(userEmail, signupId)
	}()

	sighnup_h.temp_store.SetKeyToValWIthExpiry(signupId, string(userJson), DEFAULT_SIGNUP_TIMEOUT)

	return signupId
}

func (sighnup_h *SignupHandler) HandleCompleteSignup(signupId string) (string, bool) {

	userJson := sighnup_h.temp_store.GetVal(signupId)

	if userJson == "" {
		return "", false
	}

	mapRep := map[string]any{}

	err := json.Unmarshal([]byte(userJson), &mapRep)

	if err != nil {
		return "", false
	}

	user := sighnup_h.users_store.BuildClient(mapRep)

	userId, success := sighnup_h.users_store.Save(user)

	sighnup_h.temp_store.DelKey(signupId)

	return userId, success
}

func (sighnup_h *SignupHandler) sendEmailConfirmationMsg(email, signupId string) {
	msg := ""

	emailConfirmationLink := fmt.Sprintf("<a href=%s>link</a>",
		fmt.Sprintf("%s%s%s", config.BaseApiUrl, "complete_signup/", signupId))
	if strings.Contains(config.EmailConfirmationMessage, config.LinkSubstitute) {
		msg = strings.Replace(config.EmailConfirmationMessage,
			config.LinkSubstitute,
			emailConfirmationLink, 1)
	} else {
		msg = fmt.Sprintf(`
		<div>
			<h1>Welcome</h1>
			<p1>complete signup by clicking this %s.</p1>
			<br>
			<br>
			<br>
			<p>powered by Go-Auth https://github.com/iyusuf40/goBackendUtils.</p>
		</div>
		`, emailConfirmationLink)
	}

	mail.SendMailGmail(config.GmailSource, email, config.GmailPassword,
		"Complete Signup", msg, true)
}

func MakeSignupHandler(temp_store_db, users_store_db, recordsName string) *SignupHandler {
	sighnup_h := new(SignupHandler)
	sighnup_h.temp_store = storage.GET_TempStore(config.TempStoreType, temp_store_db, recordsName)
	sighnup_h.users_store = storage.MakeUserStorage(users_store_db, recordsName)
	return sighnup_h
}
