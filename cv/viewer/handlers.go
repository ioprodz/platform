package cv_viewer

import (
	"ioprodz/common/ui"
	cv_models "ioprodz/cv/_models"
	"net/http"
)

func CreateGetCvHandler(repo cv_models.CVRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cv, _ := repo.Get("cv-id")
		ui.RenderPage(w, r, "cv/viewer/templates/simple", cv)
	}
}
