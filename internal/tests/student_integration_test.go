package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateStudentIntegration(t *testing.T) {

	router, db := setupTestRouter(t)
	defer db.Close()

	cleanDatabase(db)

	registerUser(t, router)
	loginResp := loginUser(t, router)

	body := `{
		"first_name":"Abhi",
		"last_name":"K",
		"email":"student@test.com",
		"age":20
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/students",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+loginResp.AccessToken,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", w.Code)
	}
}
