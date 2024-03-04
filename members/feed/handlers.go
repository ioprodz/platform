package members_feed

import (
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	members_models "ioprodz/members/_models"
	"net/http"
)

func CreateGetHandler(repo members_models.MembersRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		member, memberNotFound := repo.Get(user.Id)
		profileCompleted := memberNotFound == nil

		ui.RenderPage(w, r, "members/feed/index", map[string]interface{}{"Member": member, "ProfileCompleted": profileCompleted})
	}
}
