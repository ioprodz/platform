package qna_solver

import (
	"encoding/json"
	"fmt"
	"ioprodz/common/ui"
	qna_models "ioprodz/qna/_models"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

func CreateCreateAnswerHandler(qnaRepo qna_models.QNARepository, answersRepo qna_models.AnswersRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		qnaId := vars["id"]

		qna, err := qnaRepo.Get(qnaId)
		if err != nil {
			ui.Render404(w, r)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		answers := make([]int, len(qna.Questions))

		// Map values from source to destination
		for questionIndex := range qna.Questions {
			// Perform some operation on each value and append to destination
			answerValue, err := strconv.Atoi(r.Form.Get("question-" + fmt.Sprint(questionIndex) + "-choice"))
			if err != nil {
				http.Error(w, "Failed to parse form data", http.StatusBadRequest)
				return
			}
			answers[questionIndex] = answerValue
		}

		fmt.Println("choooice 0 ", r.Form.Get("question-0-choice"))

		answer := qna_models.NewAnswer(qna, answers)

		ansStr, err := json.Marshal(answer)
		fmt.Println(string(ansStr))

		answersRepo.Create(answer)

		w.Write([]byte(answer.Id))
	}
}

func CreateGetOneAnswerHandler(answersRepo qna_models.AnswersRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		answerId := vars["id"]
		answer, err := answersRepo.Get(answerId)

		if err != nil {
			ui.Render404(w, r)
			return
		}
		ui.RenderPage(w, r, "qna/solver/qna-answer", answer)
	}
}

func CreateGetOneHandler(qnaRepo qna_models.QNARepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		qnaId := vars["id"]
		qna, err := qnaRepo.Get(qnaId)
		if err != nil {
			ui.Render404(w, r)
			return
		}
		ui.RenderPage(w, r, "qna/solver/qna-form", qna)
	}
}
