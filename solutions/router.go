package solutions

import (
	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) {
	router.HandleFunc("/solutions", OverviewHandler).Methods("GET")
	router.HandleFunc("/solutions/ai-engine", AIEngineHandler).Methods("GET")
	router.HandleFunc("/solutions/chat-collaboration", ChatCollaborationHandler).Methods("GET")
	router.HandleFunc("/solutions/collaborative-editing", CollaborativeEditingHandler).Methods("GET")
	router.HandleFunc("/solutions/search-rag", SearchRAGHandler).Methods("GET")
}
