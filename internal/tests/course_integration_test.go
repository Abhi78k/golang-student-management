package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateCourseIntegration(t *testing.T) {

	router, db := setupTestRouter(t)
	defer db.Close()

	cleanDatabase(db)

	registerUser(t, router)
	loginResp := loginUser(t, router)

	body := `{
		"name":"Go Backend",
		"available_seats":10
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/courses",
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

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
