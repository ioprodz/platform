package qna_admin

import (
	"encoding/json"
	"fmt"
	"ioprodz/common/clients/openaiClient"
	"ioprodz/common/ui"

	qna_models "ioprodz/qna/_models"

	"net/http"

	"github.com/gorilla/mux"
)

func CreateListHandler(repo qna_models.QNARepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ui.RenderPage(w, r, "qna/admin/list", map[string]interface{}{"list": repo.List()})
	}
}

func CreateCreatePageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ui.RenderPage(w, r, "qna/admin/create", nil)
	}
}

type questionList struct {
	Questions []qna_models.Question `json:"questions"`
}

func CreateCreateQNAHandler(repo qna_models.QNARepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		qna := qna_models.NewQNA(r.Form.Get("name"))

		aiResponse, err := openaiClient.JsonPrompt("you are going to ask me 5 questions about '"+qna.Name+"' to assess my knowledge", "{questions:{ value:string , choices:{ value:string , isCorrect:boolean }[]}}[]}")
		if err != nil {
			fmt.Println("Error getting prompt from ai response", err)
		}
		var list *questionList
		if err := json.Unmarshal([]byte(aiResponse), &list); err != nil {
			fmt.Println("Error parsing json from ai response", err)
		}
		qna.SetQuestions(list.Questions)
		repo.Create(*qna)
		w.Write([]byte(qna.Id))
	}
}

func CreateGetOneHandler(repo qna_models.QNARepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		qnaId := vars["id"]
		qna, err := repo.Get(qnaId)
		if err != nil {
			ui.Render404(w, r)
			return
		}
		ui.RenderPage(w, r, "qna/admin/details", qna)
	}
}
