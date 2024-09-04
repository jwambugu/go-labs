package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"log"
	"strings"
)

const (
	UnusablePasswordPrefix       = "!"
	UnusablePasswordSuffixLength = 12
	PBKDF2Iterations             = 150_000
	KeyLength                    = 32
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(bytes)[:length], nil
}

func MakePassword(password string, salt string) (string, error) {
	if salt == "" {
		randomSalt, err := GenerateRandomString(16)
		if err != nil {
			return "", err
		}
		salt = randomSalt
	}

	hashedPassword := pbkdf2.Key([]byte(password), []byte(salt), PBKDF2Iterations, KeyLength, sha256.New)
	hashedPasswordStr := base64.StdEncoding.EncodeToString(hashedPassword)

	return fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", PBKDF2Iterations, salt, hashedPasswordStr), nil
}

func CheckPassword(storedHash string, password string) error {
	if strings.HasPrefix(storedHash, UnusablePasswordPrefix) {
		return errors.New("cannot use unusable password")
	}

	parts := strings.Split(storedHash, "$")
	if len(parts) != 4 {
		return errors.New("invalid hash format")
	}

	var (
		salt            = parts[2]
		storedHashValue = parts[3]

		hashedPassword        = pbkdf2.Key([]byte(password), []byte(salt), PBKDF2Iterations, KeyLength, sha256.New)
		encodedHashedPassword = base64.StdEncoding.EncodeToString(hashedPassword)
	)

	if storedHashValue != encodedHashedPassword {
		return errors.New("passwords don't mismatch")
	}

	return nil
}

func main() {
	password := "??"

	// Create a hashed password
	//hashedPassword, err := MakePassword(password, "")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Hashed Password:", hashedPassword)

	//hashedPassword := `pbkdf2_sha256$150000$CxPEOd8WoUuL$potLCiaOklwFjvtXQgmQpToPo0NpK1sSFk4YLJ2JXg0=`
	hashedPassword := `pbkdf2_sha256$150000$EnBKGlIvVx5s$tzvTB7IZoPRJiSUmCgHkS/qqVOAdaDPUb0pvzRn6N88=`

	// Simulate login attempt

	if err := CheckPassword(hashedPassword, password); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Login successful!")
}
