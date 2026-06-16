package auth

import (
	"database/sql"
	"testing"
	"time"

	"cli-auth/internal/config"
	"cli-auth/internal/database"

	"github.com/pquerna/otp/totp"
)

func setupTestDB(t *testing.T) (*sql.DB, *Repository) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	err = database.CreateSchema(db)
	if err != nil {
		db.Close()
		t.Fatalf("failed to create schema: %v", err)
	}

	repo := NewRepository(db)
	return db, repo
}

func TestPasswordHashing(t *testing.T) {
	password := "my-secure-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty password hash")
	}

	err = VerifyPassword(hash, password)
	if err != nil {
		t.Errorf("VerifyPassword failed for correct password: %v", err)
	}

	err = VerifyPassword(hash, "wrong-password")
	if err == nil {
		t.Error("expected VerifyPassword to fail for incorrect password")
	}
}

func TestUserRegistration(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	// Test successful registration
	err := Register(repo, "testuser", "securepass123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Verify user is in db
	user, err := repo.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("expected username to be 'testuser', got '%s'", user.Username)
	}

	if user.FailedAttempts != 0 {
		t.Errorf("expected failed attempts to be 0, got %d", user.FailedAttempts)
	}

	// Test empty username validation
	err = Register(repo, "", "pass")
	if err == nil {
		t.Error("expected error when registering empty username")
	}

	// Test empty password validation
	err = Register(repo, "user2", "")
	if err == nil {
		t.Error("expected error when registering empty password")
	}

	// Test duplicate username registration
	err = Register(repo, "testuser", "anotherpass")
	if err == nil {
		t.Error("expected error when registering duplicate username")
	} else if err.Error() != "username already exists" {
		t.Errorf("expected error 'username already exists', got '%s'", err.Error())
	}
}

func TestAuthenticationAndLockout(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	username := "lockuser"
	password := "mypassword"

	// Override config for faster testing
	origMaxAttempts := config.App.MaxFailedAttempt
	origLockDuration := config.App.LockDuration
	defer func() {
		config.App.MaxFailedAttempt = origMaxAttempts
		config.App.LockDuration = origLockDuration
	}()

	config.App.MaxFailedAttempt = 3
	config.App.LockDuration = 200 * time.Millisecond

	// Register user
	err := Register(repo, username, password)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Test successful authentication
	user, err := AuthenticatePassword(repo, username, password)
	if err != nil {
		t.Fatalf("AuthenticatePassword failed for correct credentials: %v", err)
	}

	if user.Username != username {
		t.Errorf("expected username '%s', got '%s'", username, user.Username)
	}

	// Test failed authentication (Attempt 1)
	_, err = AuthenticatePassword(repo, username, "wrong-pass")
	if err == nil {
		t.Fatal("expected authentication error")
	}

	user, _ = repo.GetUserByUsername(username)
	if user.FailedAttempts != 1 {
		t.Errorf("expected 1 failed attempt, got %d", user.FailedAttempts)
	}

	// Failed authentication (Attempt 2)
	_, err = AuthenticatePassword(repo, username, "wrong-pass")
	if err == nil {
		t.Fatal("expected authentication error")
	}

	user, _ = repo.GetUserByUsername(username)
	if user.FailedAttempts != 2 {
		t.Errorf("expected 2 failed attempts, got %d", user.FailedAttempts)
	}

	// Failed authentication (Attempt 3 - should trigger lockout)
	_, err = AuthenticatePassword(repo, username, "wrong-pass")
	if err == nil {
		t.Fatal("expected authentication error")
	}

	user, _ = repo.GetUserByUsername(username)
	if user.FailedAttempts != 3 {
		t.Errorf("expected 3 failed attempts, got %d", user.FailedAttempts)
	}

	if user.LockedUntil == nil {
		t.Fatal("expected user to be locked out (LockedUntil is nil)")
	}

	// Try to authenticate while locked
	_, err = AuthenticatePassword(repo, username, password)
	if err == nil {
		t.Fatal("expected login to be blocked while account is locked")
	}

	// Wait for lockout to expire
	time.Sleep(250 * time.Millisecond)

	// After lockout expired, wrong password should start a new failed attempt count (1), not lock immediately again
	_, err = AuthenticatePassword(repo, username, "wrong-pass-again")
	if err == nil {
		t.Fatal("expected authentication error")
	}

	user, _ = repo.GetUserByUsername(username)
	if user.FailedAttempts != 1 {
		t.Errorf("expected failed attempts to reset and become 1 after lock expired, got %d", user.FailedAttempts)
	}

	// Log in successfully with correct password now
	user, err = AuthenticatePassword(repo, username, password)
	if err != nil {
		t.Fatalf("expected login to succeed: %v", err)
	}

	err = CompleteLogin(repo, user.ID)
	if err != nil {
		t.Fatalf("CompleteLogin failed: %v", err)
	}

	// Verify failed attempts count reset
	user, _ = repo.GetUserByUsername(username)
	if user.FailedAttempts != 0 {
		t.Errorf("expected failed attempts reset to 0, got %d", user.FailedAttempts)
	}

	if user.LockedUntil != nil {
		t.Errorf("expected LockedUntil to be cleared, got %s", user.LockedUntil)
	}

	if user.LastLogin == nil {
		t.Error("expected LastLogin to be updated")
	}
}

func TestTOTPFunctionality(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	username := "totpuser"
	err := Register(repo, username, "password")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	user, _ := repo.GetUserByUsername(username)

	// Generate TOTP key
	key, err := GenerateTOTPKey(username)
	if err != nil {
		t.Fatalf("GenerateTOTPKey failed: %v", err)
	}

	if key.Secret() == "" {
		t.Fatal("expected key secret to be non-empty")
	}

	// Verify validation works with correct TOTP code
	code, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		t.Fatalf("failed to generate code from TOTP key: %v", err)
	}

	if !ValidateTOTP(code, key.Secret()) {
		t.Error("expected TOTP code to be valid")
	}

	if ValidateTOTP("000000", key.Secret()) {
		t.Error("expected invalid TOTP code to fail validation")
	}

	// Test Enabling 2FA
	err = repo.EnableMFA(user.ID, key.Secret())
	if err != nil {
		t.Fatalf("EnableMFA failed: %v", err)
	}

	user, _ = repo.GetUserByUsername(username)
	if !user.MFAEnabled {
		t.Error("expected user.MFAEnabled to be true")
	}

	if user.TOTPSecret == nil || *user.TOTPSecret != key.Secret() {
		t.Errorf("expected totp secret to be saved, got %v", user.TOTPSecret)
	}

	// Test Disabling 2FA
	err = repo.DisableMFA(user.ID)
	if err != nil {
		t.Fatalf("DisableMFA failed: %v", err)
	}

	user, _ = repo.GetUserByUsername(username)
	if user.MFAEnabled {
		t.Error("expected user.MFAEnabled to be false")
	}

	if user.TOTPSecret != nil {
		t.Errorf("expected totp secret to be cleared (nil), got %v", user.TOTPSecret)
	}
}
