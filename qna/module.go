package qna

import (
	qna_infra "ioprodz/qna/_infra"
	qna_admin "ioprodz/qna/admin"
	qna_solver "ioprodz/qna/solver"

	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) {

	qnaRepo := qna_infra.CreateMemoryQNARepo()
	answersRepo := qna_infra.CreateMemoryAnswerRepo()

	// pages
	router.HandleFunc("/admin/qna", qna_admin.CreateListHandler(qnaRepo)).Methods("GET")
	router.HandleFunc("/admin/qna/create-new", qna_admin.CreateCreatePageHandler()).Methods("GET")
	router.HandleFunc("/admin/qna/{id}", qna_admin.CreateGetOneHandler(qnaRepo)).Methods("GET")

	router.HandleFunc("/qna/{id}/answers", qna_solver.CreateCreateAnswerHandler(qnaRepo, answersRepo)).Methods("POST")
	router.HandleFunc("/qna/{id}", qna_solver.CreateGetOneHandler(qnaRepo)).Methods("GET")
	router.HandleFunc("/qna-answers/{id}", qna_solver.CreateGetOneAnswerHandler(answersRepo)).Methods("GET")

	// api
	router.HandleFunc("/api/admin/qna", qna_admin.CreateCreateQNAHandler(qnaRepo)).Methods("POST")

}
