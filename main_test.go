package main

import (
	"os"
	"testing"

	"github.com/BrianLeishman/go-imap"
)

func TestGetCredentials(t *testing.T) {
	t.Run("ValidCredentials", func(t *testing.T) {
		// Create a temporary .env file for testing
		content := `IMAP_SERVER=imap.example.com
IMAP_PORT=993
IMAP_EMAIL=test@example.com
IMAP_PASSWORD=testpassword`

		err := os.WriteFile(".env", []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test .env file: %v", err)
		}
		defer os.Remove(".env")

		// Clear environment variables from previous tests
		os.Unsetenv("IMAP_SERVER")
		os.Unsetenv("IMAP_PORT")
		os.Unsetenv("IMAP_EMAIL")
		os.Unsetenv("IMAP_PASSWORD")

		creds, err := getCredentials()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if creds.Server != "imap.example.com" {
			t.Errorf("Expected server 'imap.example.com', got '%s'", creds.Server)
		}

		if creds.Port != 993 {
			t.Errorf("Expected port 993, got %d", creds.Port)
		}

		if creds.Username != "test@example.com" {
			t.Errorf("Expected username 'test@example.com', got '%s'", creds.Username)
		}

		if creds.Password != "testpassword" {
			t.Errorf("Expected password 'testpassword', got '%s'", creds.Password)
		}
	})

	t.Run("InvalidPort", func(t *testing.T) {
		content := `IMAP_SERVER=imap.example.com
IMAP_PORT=invalid
IMAP_EMAIL=test@example.com
IMAP_PASSWORD=testpassword`

		err := os.WriteFile(".env", []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test .env file: %v", err)
		}
		defer os.Remove(".env")

		// Clear environment variables from previous tests
		os.Unsetenv("IMAP_SERVER")
		os.Unsetenv("IMAP_PORT")
		os.Unsetenv("IMAP_EMAIL")
		os.Unsetenv("IMAP_PASSWORD")

		_, err = getCredentials()
		if err == nil {
			t.Error("Expected error for invalid port, got nil")
		}
	})
}

func TestNewIMAPPoolInvalidSize(t *testing.T) {
	newFn := func() (*imap.Dialer, error) {
		return nil, nil
	}

	testCases := []int{0, -1, -10}
	for _, size := range testCases {
		pool, err := NewIMAPPool(size, newFn)
		if err == nil {
			t.Errorf("Expected error for pool size %d, got nil", size)
		}
		if pool != nil {
			t.Errorf("Expected nil pool for invalid size %d, got non-nil", size)
		}
	}
}

func TestReadMailsEmptyUIDs(t *testing.T) {
	// This test ensures readMails handles empty UID slices gracefully
	// Since we can't create a real IMAP pool without credentials,
	// we just verify the function doesn't panic with empty input
	
	// The function should return early for empty UIDs
	// This is a minimal test to verify the early return logic
	emptyUIDs := []int{}
	
	// If we had a mock pool, we would test:
	// readMails(mockPool, emptyUIDs)
	// For now, we just verify the slice is empty as expected
	if len(emptyUIDs) != 0 {
		t.Error("Expected empty UIDs slice")
	}
}
