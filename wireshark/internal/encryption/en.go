package securevault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	currentVersion = 1

	saltSize  = 16
	nonceSize = 12
	keySize   = 32 // AES-256
)

type Argon2Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

type Vault struct {
	password string
	params   Argon2Params
}

// New creates a Vault with default Argon2 params.
func New(password string) *Vault {
	return &Vault{
		password: password,
		params: Argon2Params{
			Time:    3,
			Memory:  64 * 1024,
			Threads: 4,
		},
	}
}

// NewWithParams allows custom Argon2 params.
func NewWithParams(password string, params Argon2Params) *Vault {
	return &Vault{
		password: password,
		params:   params,
	}
}

// Encrypt encrypts plaintext with optional AAD, returns base64(version + header + ciphertext)
func (v *Vault) Encrypt(plaintext, aad []byte) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(v.password), salt, v.params.Time, v.params.Memory, v.params.Threads, keySize)
	defer zeroBytes(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, aad)

	// Construct binary header
	var buf bytes.Buffer
	buf.WriteByte(currentVersion)
	buf.Write(salt)
	buf.Write(nonce)
	_ = binary.Write(&buf, binary.BigEndian, v.params.Time)
	_ = binary.Write(&buf, binary.BigEndian, v.params.Memory)
	buf.WriteByte(v.params.Threads)
	buf.Write(ciphertext)

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Decrypt parses base64(input), extracts header, re-derives key, and decrypts using optional AAD.
func (v *Vault) Decrypt(encoded string, aad []byte) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	if len(raw) < 1+saltSize+nonceSize+4+4+1 {
		return nil, errors.New("ciphertext too short")
	}

	reader := bytes.NewReader(raw)

	var version byte
	if err := binary.Read(reader, binary.BigEndian, &version); err != nil {
		return nil, err
	}
	if version != currentVersion {
		return nil, errors.New("unsupported version")
	}

	salt := make([]byte, saltSize)
	_, _ = io.ReadFull(reader, salt)
	nonce := make([]byte, nonceSize)
	_, _ = io.ReadFull(reader, nonce)

	var time uint32
	var memory uint32
	var threads uint8

	_ = binary.Read(reader, binary.BigEndian, &time)
	_ = binary.Read(reader, binary.BigEndian, &memory)
	_ = binary.Read(reader, binary.BigEndian, &threads)

	ciphertext, _ := io.ReadAll(reader)

	key := argon2.IDKey([]byte(v.password), salt, time, memory, threads, keySize)
	defer zeroBytes(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, errors.New("decryption failed: invalid password or corrupted data")
	}
	return plaintext, nil
}

// ConstantTimeEqual compares two byte slices securely.
func ConstantTimeEqual(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// zeroBytes clears memory.
func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
