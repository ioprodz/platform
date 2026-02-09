package consulting

import (
	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) {
	router.HandleFunc("/consulting", OverviewHandler).Methods("GET")
	router.HandleFunc("/consulting/it-strategy", ITStrategyHandler).Methods("GET")
	router.HandleFunc("/consulting/coaching", CoachingHandler).Methods("GET")
}
