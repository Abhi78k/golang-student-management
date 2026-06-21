package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project-1/internal/dto"

	"github.com/gin-gonic/gin"
)

func registerUser(
	t *testing.T,
	router *gin.Engine,
) {

	body := `{
		"username":"abhi",
		"email":"abhi@test.com",
		"password":"password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf(
			"register failed: %s",
			w.Body.String(),
		)
	}
}

func loginUser(
	t *testing.T,
	router *gin.Engine,
) dto.LoginResponse {

	body := `{
		"email":"abhi@test.com",
		"password":"password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf(
			"login failed: %s",
			w.Body.String(),
		)
	}

	var response dto.LoginResponse

	err := json.Unmarshal(
		w.Body.Bytes(),
		&response,
	)

	if err != nil {
		t.Fatal(err)
	}

	return response
}
