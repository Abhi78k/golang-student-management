package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"project-1/internal/dto"
	"strings"
	"testing"
)

func TestRegisterIntegration(
	t *testing.T,
) {
	router, db := setupTestRouter(t)

	defer db.Close()

	cleanDatabase(db)

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

	router.ServeHTTP(
		w,
		req,
	)

	if w.Code != http.StatusCreated {
		t.Fatalf(
			"expected 201 got %d",
			w.Code,
		)
	}
}

func TestLoginIntegration(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()

	cleanDatabase(db)

	registerBody := `{
		"username":"abhi",
		"email":"abhi@test.com",
		"password":"password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(registerBody),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	loginBody := `{
		"email":"abhi@test.com",
		"password":"password123"
	}`

	req = httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		strings.NewReader(loginBody),
	)

	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	var response dto.LoginResponse

	json.Unmarshal(
		w.Body.Bytes(),
		&response,
	)

	if response.AccessToken == "" {
		t.Fatal("access token missing")
	}

	if response.RefreshToken == "" {
		t.Fatal("refresh token missing")
	}
}

func TestGetMeIntegration(t *testing.T) {

	router, db := setupTestRouter(t)
	defer db.Close()

	cleanDatabase(db)

	registerUser(t, router)

	loginResp := loginUser(t, router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/auth/me",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+loginResp.AccessToken,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestRefreshIntegration(t *testing.T) {

	router, db := setupTestRouter(t)
	defer db.Close()

	cleanDatabase(db)

	registerUser(t, router)

	loginResp := loginUser(t, router)

	body := fmt.Sprintf(
		`{"refresh_token":"%s"}`,
		loginResp.RefreshToken,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/refresh",
		strings.NewReader(body),
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

func TestLogoutIntegration(t *testing.T) {

	router, db := setupTestRouter(t)
	defer db.Close()

	cleanDatabase(db)

	registerUser(t, router)

	loginResp := loginUser(t, router)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/logout",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+loginResp.AccessToken,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
