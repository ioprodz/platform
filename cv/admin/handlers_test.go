package cv_admin_test

import (
	cv_admin "ioprodz/cv/admin"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCreateCvHandler(t *testing.T) {

	t.Run("should reply", func(t *testing.T) {

		// arrange
		w := httptest.NewRecorder()
		handler := cv_admin.CreateCreateCvHandler()

		// act
		handler(w, httptest.NewRequest("GET", "/", nil))

		// assert
		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.String(), "test")
	})
}
