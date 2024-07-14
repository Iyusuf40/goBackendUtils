package auth

import (
	"fmt"
	"strings"

	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/Iyusuf40/goBackendUtils/mail"
	"github.com/Iyusuf40/goBackendUtils/models"
	"github.com/Iyusuf40/goBackendUtils/storage"
	"github.com/google/uuid"
)

type AuthHandler struct {
	temp_store  storage.TempStore
	users_store storage.Storage[models.User]
}

const DEFAULT_SESSION_TIMEOUT = 86400.0

func (auth_h *AuthHandler) HandleLogin(email, password string) string {
	retrievedUsers := auth_h.users_store.GetByField("email", email)

	if len(retrievedUsers) == 0 {
		return ""
	}

	user := retrievedUsers[0]
	if !user.IsCorrectPassword(password) {
		return ""
	}

	sessionId := uuid.NewString()
	userId := auth_h.users_store.GetIdByField("email", email)

	auth_h.temp_store.SetKeyToValWIthExpiry(sessionId, userId, DEFAULT_SESSION_TIMEOUT)

	return sessionId
}

func (auth_h *AuthHandler) HandleLogout(sessionId string) {
	auth_h.temp_store.DelKey(sessionId)
}

func (auth_h *AuthHandler) IsLoggedIn(sessionId string) bool {
	return auth_h.temp_store.GetVal(sessionId) != ""
}

func (auth_h *AuthHandler) HandleForgotPassword(email string) string {
	userId := auth_h.users_store.GetIdByField("email", email)
	if userId == "" {
		fmt.Println("HandleForgotPassword: user does not exist")
		return ""
	}

	passwordResetToken := uuid.New().String()

	auth_h.sendPasswordResetEmail(email, passwordResetToken)

	auth_h.temp_store.SetKeyToValWIthExpiry(passwordResetToken, userId, 3600)

	return passwordResetToken
}

func (auth_h *AuthHandler) sendPasswordResetEmail(email, passwordResetToken string) {
	msg := ""

	passwordResetLink := fmt.Sprintf("<a href=%s>link</a>",
		fmt.Sprintf("%s%s%s", config.BaseAuthUrl, "reset_password/", passwordResetToken))
	if strings.Contains(config.PasswordResetMessage, config.LinkSubstitute) {
		msg = strings.Replace(config.PasswordResetMessage,
			config.LinkSubstitute,
			passwordResetLink, 1)
	} else {
		msg = fmt.Sprintf(`
		<div>
			<h1>Reset Password</h1>
			<p1>Reset Password by clicking this %s.</p1>
			<br>
			<br>
			<br>
			<p>powered by Go-Auth https://github.com/iyusuf40/goBackendUtils.</p>
		</div>
		`, passwordResetLink)
	}

	mail.SendMailGmail(config.GmailSource, email, config.GmailPassword,
		"Reset Password", msg, true)
}

func (auth_h *AuthHandler) HandleUpdatePassword(passwordResetToken, newPassword string) bool {

	if newPassword == "" {
		fmt.Println("password cannot be empty")
		return false
	}

	userId := auth_h.temp_store.GetVal(passwordResetToken)
	if userId == "" {
		fmt.Println("passwordResetToken has no value in store")
		return false
	}

	_, err := auth_h.users_store.Get(userId)

	if err != nil {
		fmt.Println("user with userId:" + userId + " does not exist")
		return false
	}

	return auth_h.users_store.Update(userId,
		storage.UpdateDesc{Field: "password", Value: newPassword})
}

func (auth_h *AuthHandler) ExtendSession(sessionId string, duration float64) {
	auth_h.temp_store.ChangeKeyEpiry(sessionId, duration)
}

func MakeAuthHandler(temp_store_db, users_store_db, recordsName string) *AuthHandler {
	auth_h := new(AuthHandler)
	auth_h.temp_store = storage.GET_TempStore(config.TempStoreType, temp_store_db, recordsName)
	auth_h.users_store = storage.MakeUserStorage(users_store_db, recordsName)
	return auth_h
}
