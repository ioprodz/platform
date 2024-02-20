package qna

import (
	qna_admin "ioprodz/qna/admin"
	qna_solver "ioprodz/qna/solver"

	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) {

	// pages
	router.HandleFunc("/admin/qna", qna_admin.ListHandler).Methods("GET")
	router.HandleFunc("/admin/qna/create-new", qna_admin.CreatePageHandler).Methods("GET")
	router.HandleFunc("/admin/qna/{id}", qna_admin.GetOneHandler).Methods("GET")

	router.HandleFunc("/qna/{id}/answers", qna_solver.CreateAnswerHandler).Methods("POST")
	router.HandleFunc("/qna/{id}", qna_solver.GetOneHandler).Methods("GET")
	router.HandleFunc("/qna-answers/{id}", qna_solver.GetOneAnswerHandler).Methods("GET")

	// api
	router.HandleFunc("/api/admin/qna", qna_admin.CreateQNAHandler).Methods("POST")

}
