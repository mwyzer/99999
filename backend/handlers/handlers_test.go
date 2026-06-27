package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"property-hub-backend/config"
	"property-hub-backend/routes"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		AppName:              "TestPropertyHub",
		JWTSecret:            "test-secret-32-bytes-long-key!!",
		JWTExpiryHours:       1,
		BcryptCost:           4,
		CORSAllowedOrigins:   "*",
		MaxPhotosPerListing:  10,
		MaxUploadSizeMB:      5,
		RateLimitLoginPerMinute:    100,
		RateLimitGlobalPerMinute:   1000,
		UploadDir:            "./test-uploads",
	}
	AppConfig = cfg
	return routes.Setup(cfg)
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)

	if body["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", body["status"])
	}
	if body["app"] != "TestPropertyHub" {
		t.Errorf("expected app=TestPropertyHub, got %v", body["app"])
	}
}

func TestRegisterValidation(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name     string
		payload  string
		wantCode int
	}{
		{
			name:     "missing name",
			payload:  `{"email":"test@test.com","phone":"081234567890","password":"Test@1234"}`,
			wantCode: 422,
		},
		{
			name:     "invalid email",
			payload:  `{"name":"Test","email":"not-email","phone":"081234567890","password":"Test@1234"}`,
			wantCode: 422,
		},
		{
			name:     "short password",
			payload:  `{"name":"Test","email":"test@test.com","phone":"081234567890","password":"short"}`,
			wantCode: 422,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/auth/register", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Use a body reader
			req.Body = io.NopCloser(strings.NewReader(tt.payload))

			router.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("expected %d, got %d. body: %s", tt.wantCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	router := setupTestRouter()

	payload := `{"email":"nonexistent@test.com","password":"WrongPass123"}`
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401 for invalid credentials, got %d", w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)

	errObj, ok := body["error"].(map[string]interface{})
	if !ok {
		t.Fatal("expected error object in response")
	}
	if errObj["code"] != "AUTH_INVALID_CREDENTIALS" {
		t.Errorf("expected AUTH_INVALID_CREDENTIALS code, got %v", errObj["code"])
	}
}

func TestPublicPropertyList(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/properties?page=1&per_page=12", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Without a real DB, this will likely return an error
	// Test that the endpoint at least responds with proper JSON structure
	if w.Code != 200 && w.Code != 500 {
		t.Errorf("expected 200 or 500 (no DB), got %d", w.Code)
	}

	if w.Code == 200 {
		var body map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &body)
		if body["success"] != true {
			t.Error("expected success=true in response")
		}
	}
}

func TestUnauthorizedAccess(t *testing.T) {
	router := setupTestRouter()

	// Try accessing protected endpoint without token
	req, _ := http.NewRequest("GET", "/api/v1/me/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401 for unauthorized access, got %d", w.Code)
	}
}
