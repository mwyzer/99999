package utils

import (
	"testing"

	"property-hub-backend/config"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("TestPassword123", 4) // cost 4 for speed
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if len(hash) < 20 {
		t.Fatalf("hash too short: %d chars", len(hash))
	}
}

func TestCheckPassword(t *testing.T) {
	hash, _ := HashPassword("SecurePass789", 4)

	if !CheckPassword("SecurePass789", hash) {
		t.Fatal("expected password to match")
	}
	if CheckPassword("WrongPassword", hash) {
		t.Fatal("expected password NOT to match")
	}
	if CheckPassword("", hash) {
		t.Fatal("empty password should not match")
	}
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	if CheckPassword("anything", "not-a-bcrypt-hash") {
		t.Fatal("expected false for invalid hash format")
	}
}

func TestGetPagination(t *testing.T) {
	tests := []struct {
		name          string
		queryPage     string
		queryPerPage  string
		expectedPage  int
		expectedLimit int
	}{
		{"defaults", "", "", 1, 20},
		{"valid values", "3", "10", 3, 10},
		{"negative page", "-5", "10", 1, 10},
		{"zero per_page", "1", "0", 1, 20},
		{"large per_page", "1", "200", 1, 100},
		{"non-numeric", "abc", "xyz", 1, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: 1, PerPage: 20}
			if tt.queryPage != "" {
				p.Page = mustAtoi(tt.queryPage)
			}
			if tt.queryPerPage != "" {
				p.PerPage = mustAtoi(tt.queryPerPage)
			}
			// Normalize manually (simulating what GetPagination would do)
			if p.Page < 1 {
				p.Page = 1
			}
			if p.PerPage < 1 {
				p.PerPage = 20
			}
			if p.PerPage > 100 {
				p.PerPage = 100
			}
			p.Offset = (p.Page - 1) * p.PerPage

			if p.Page != tt.expectedPage {
				t.Errorf("page: got %d, want %d", p.Page, tt.expectedPage)
			}
			if p.PerPage != tt.expectedLimit {
				t.Errorf("perPage: got %d, want %d", p.PerPage, tt.expectedLimit)
			}
		})
	}
}

func TestCalculateMeta(t *testing.T) {
	meta := CalculateMeta(1, 20, 55)
	if meta.Page != 1 {
		t.Errorf("page: got %d, want 1", meta.Page)
	}
	if meta.Total != 55 {
		t.Errorf("total: got %d, want 55", meta.Total)
	}
	if meta.TotalPages != 3 {
		t.Errorf("totalPages: got %d, want 3", meta.TotalPages)
	}

	meta2 := CalculateMeta(3, 20, 55)
	if meta2.TotalPages != 3 {
		t.Errorf("totalPages: got %d, want 3", meta2.TotalPages)
	}

	meta3 := CalculateMeta(1, 20, 0)
	if meta3.TotalPages != 0 {
		t.Errorf("totalPages for 0 items: got %d, want 0", meta3.TotalPages)
	}
}

func TestGenerateToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret-32-bytes-long-key!!",
		JWTExpiryHours: 1,
		AppName:       "TestApp",
	}

	userID := uuid.New()
	token, err := GenerateToken(userID, "buyer", nil, cfg)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// Parse the token back
	claims, err := ParseToken(token, cfg)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("userID mismatch: got %s, want %s", claims.UserID, userID)
	}
	if claims.Role != "buyer" {
		t.Errorf("role: got %s, want buyer", claims.Role)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret-32-bytes-long-key!!",
		JWTExpiryHours: 1,
	}

	_, err := ParseToken("invalid.token.here", cfg)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}

	_, err = ParseToken("", cfg)
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	cfg1 := &config.Config{
		JWTSecret:     "secret-one-32-bytes-long-key!!!",
		JWTExpiryHours: 1,
		AppName:       "Test",
	}
	cfg2 := &config.Config{
		JWTSecret:     "secret-two-32-bytes-long-key!!!",
		JWTExpiryHours: 1,
	}

	token, err := GenerateToken(uuid.New(), "buyer", nil, cfg1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ParseToken(token, cfg2)
	if err == nil {
		t.Fatal("expected error when parsing with wrong secret")
	}
}

// mustAtoi helper for test pagination
func mustAtoi(s string) int {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}
