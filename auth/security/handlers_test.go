package auth_security

import (
	"context"
	auth_infra "ioprodz/auth/_infra"
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/policies"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func revokeSession(repo auth_models.SessionRepository, userId string, sessionId string) *httptest.ResponseRecorder {
	request := httptest.NewRequest("DELETE", "/security/sessions/"+sessionId, nil)
	request = mux.SetURLVars(request, map[string]string{"id": sessionId})
	request = request.WithContext(context.WithValue(request.Context(), policies.CurrentUserCtxKey, policies.CurrentUser{Id: userId}))
	recorder := httptest.NewRecorder()
	CreateRevokeSessionHandler(repo)(recorder, request)
	return recorder
}

func TestRevokeSessionDisconnectsTheUsersOwnSession(t *testing.T) {
	repo := auth_infra.CreateMemorySessionRepo()
	session := auth_models.NewSession("ana", "", "hash-ana", "", "Ana", "ana@example.com", "github")
	repo.Create(session)

	recorder := revokeSession(repo, "ana", session.Id)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", recorder.Code)
	}
	if _, err := repo.Get(session.Id); err == nil {
		t.Error("session not deleted")
	}
}

func TestRevokeSessionRefusesAnotherUsersSession(t *testing.T) {
	repo := auth_infra.CreateMemorySessionRepo()
	session := auth_models.NewSession("ana", "", "hash-ana", "", "Ana", "ana@example.com", "github")
	repo.Create(session)

	recorder := revokeSession(repo, "bob", session.Id)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", recorder.Code)
	}
	if _, err := repo.Get(session.Id); err != nil {
		t.Error("another user's session was deleted")
	}
}
