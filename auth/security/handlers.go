package auth_security

import (
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mileusna/useragent"
	"github.com/xeonx/timeago"
)

type SessionView struct {
	Id        string
	Title     string
	IsCurrent bool
	CreatedOn string
	LastUsed  string
}

func CreateSecurityPageHandler(sessionRepo auth_models.SessionRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		sessions := sessionRepo.GetByAccountId(user.Id)

		sessionViewList := []SessionView{}
		for _, session := range sessions {
			ua := useragent.Parse(session.UaString)

			createdOn, _ := time.Parse(time.RFC3339, session.FirstCreatedAt)
			lastUsed, _ := time.Parse(time.RFC3339, session.LastUsedAt)

			sessionViewList = append(sessionViewList, SessionView{
				Id:        session.Id,
				Title:     ua.Name + " on " + ua.Device + " " + ua.OS + "",
				CreatedOn: timeago.English.Format(createdOn),
				LastUsed:  timeago.English.Format(lastUsed),
				IsCurrent: session.Id == user.SessionId,
			})
		}

		ui.RenderPage(w, r, "auth/security/settings", sessionViewList)
	}
}

func CreateRevokeSessionHandler(sessionRepo auth_models.SessionRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["id"]

		sessionRepo.Delete(sessionId)
		w.Write([]byte("Disconnected Successfully"))
	}
}
