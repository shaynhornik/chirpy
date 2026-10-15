package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("04234")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}
	if hash == "04234" {
		t.Errorf("hash should not be the same as the password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, err := HashPassword("04234")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	match, err := CheckPasswordHash("04234", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error: %v", err)
	}
	if !match {
		t.Errorf("right password should match")
	}

	match, err = CheckPasswordHash("wrong", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error: %v", err)
	}
	if match {
		t.Errorf("wrong password should not match")
	}

	_, err = CheckPasswordHash("04234", "unset")
	if err == nil {
		t.Errorf("a bad hash should return an error")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	gotID, err := ValidateJWT(token, "secret")
	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}
	if gotID != userID {
		t.Errorf("got user id %v, want %v", gotID, userID)
	}
}

func TestExpiredJWT(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "secret", -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	_, err = ValidateJWT(token, "secret")
	if err == nil {
		t.Errorf("expired token should be rejected")
	}
}

func TestWrongSecretJWT(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Errorf("token signed with a different secret should be rejected")
	}
}

func TestGarbageJWT(t *testing.T) {
	_, err := ValidateJWT("not.a.jwt", "secret")
	if err == nil {
		t.Errorf("garbage token should be rejected")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer abc.def.ghi")
	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned an error: %v", err)
	}
	if token != "abc.def.ghi" {
		t.Errorf("got token %q, want %q", token, "abc.def.ghi")
	}

	headers.Set("Authorization", "Bearer    abc.def.ghi   ")
	token, err = GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned an error: %v", err)
	}
	if token != "abc.def.ghi" {
		t.Errorf("extra whitespace: got token %q, want %q", token, "abc.def.ghi")
	}
}

func TestGetBearerTokenMissing(t *testing.T) {
	_, err := GetBearerToken(http.Header{})
	if err == nil {
		t.Errorf("missing header should return an error")
	}

	headers := http.Header{}
	headers.Set("Authorization", "Basic abc123")
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Errorf("non-bearer header should return an error")
	}

	headers.Set("Authorization", "Bearer ")
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Errorf("empty bearer token should return an error")
	}
}

func TestMakeRefreshToken(t *testing.T) {
	token1 := MakeRefreshToken()
	token2 := MakeRefreshToken()
	if len(token1) != 64 {
		t.Errorf("refresh token should be 64 hex characters, got %d", len(token1))
	}
	if token1 == token2 {
		t.Errorf("two refresh tokens should not be the same")
	}
}

func TestGetAPIKey(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "ApiKey f271c81ff7084ee5b99a5091b42d486e")
	key, err := GetAPIKey(headers)
	if err != nil {
		t.Fatalf("GetAPIKey returned an error: %v", err)
	}
	if key != "f271c81ff7084ee5b99a5091b42d486e" {
		t.Errorf("got key %q", key)
	}

	_, err = GetAPIKey(http.Header{})
	if err == nil {
		t.Errorf("missing header should return an error")
	}

	headers.Set("Authorization", "Bearer f271c81ff7084ee5b99a5091b42d486e")
	_, err = GetAPIKey(headers)
	if err == nil {
		t.Errorf("bearer header should return an error")
	}

	headers.Set("Authorization", "ApiKey ")
	_, err = GetAPIKey(headers)
	if err == nil {
		t.Errorf("empty api key should return an error")
	}
}
