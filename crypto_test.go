package wecombot

import (
	"testing"
)

func TestCrypto_EncryptDecrypt(t *testing.T) {
	config := Config{
		Token:     "test_token",
		AESKey:    "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG",
		ReceiveID: "",
	}

	crypto, err := NewCrypto(config)
	if err != nil {
		t.Fatalf("NewCrypto failed: %v", err)
	}

	testMsg := []byte(`{"msgid":"test123","msgtype":"text"}`)

	// Encrypt
	encrypted, err := crypto.Encrypt(testMsg)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := crypto.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(testMsg) {
		t.Errorf("Decrypted message mismatch: got %s, want %s", string(decrypted), string(testMsg))
	}
}

func TestCrypto_VerifySignature(t *testing.T) {
	config := Config{
		Token:     "test_token",
		AESKey:    "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG",
		ReceiveID: "",
	}

	crypto, err := NewCrypto(config)
	if err != nil {
		t.Fatalf("NewCrypto failed: %v", err)
	}

	timestamp := "1234567890"
	nonce := "random_nonce"
	encrypt := "encrypted_data"

	// Calculate expected signature
	expected := crypto.calculateSignature(timestamp, nonce, encrypt)

	// Verify with correct signature
	err = crypto.VerifySignature(expected, timestamp, nonce, encrypt)
	if err != nil {
		t.Errorf("VerifySignature failed with correct signature: %v", err)
	}

	// Verify with wrong signature
	err = crypto.VerifySignature("wrong_signature", timestamp, nonce, encrypt)
	if err != ErrBadSignature {
		t.Errorf("VerifySignature should fail with wrong signature, got: %v", err)
	}
}

func TestPKCS7Padding(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		blockSize int
	}{
		{"empty", []byte{}, 32},
		{"partial block", []byte("hello"), 32},
		{"full block", []byte("12345678901234567890123456789012"), 32},
		{"multiple blocks", []byte("1234567890123456789012345678901212345"), 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			padded := pkcs7Pad(tt.data, tt.blockSize)
			if len(padded)%tt.blockSize != 0 {
				t.Errorf("Padded data length %d is not multiple of block size %d", len(padded), tt.blockSize)
			}

			unpadded, err := pkcs7Unpad(padded, tt.blockSize)
			if err != nil {
				t.Fatalf("pkcs7Unpad failed: %v", err)
			}

			if string(unpadded) != string(tt.data) {
				t.Errorf("Unpadded data mismatch: got %s, want %s", string(unpadded), string(tt.data))
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr error
	}{
		{
			name: "valid config",
			config: Config{
				Token:     "token",
				AESKey:    "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG",
				ReceiveID: "",
			},
			wantErr: nil,
		},
		{
			name: "empty token",
			config: Config{
				Token:     "",
				AESKey:    "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG",
				ReceiveID: "",
			},
			wantErr: ErrInvalidConfig,
		},
		{
			name: "invalid AES key length",
			config: Config{
				Token:     "token",
				AESKey:    "short",
				ReceiveID: "",
			},
			wantErr: ErrInvalidAESKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
