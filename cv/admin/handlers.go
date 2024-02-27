package cv_admin

import "net/http"

func CreateCreateCvHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}
}
