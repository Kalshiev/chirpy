package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateJWT(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	testUser := uuid.New()

	if err != nil {
		t.Fatal(err)
	}
	token, err := MakeJWT(testUser, "AllYourBase", dur)
	if err != nil {
		t.Fatal(err)
	}

	returnedUser, err := ValidateJWT(token, "AllYourBase")
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(testUser.String(), returnedUser.String()) != 0 {
		t.Error("Mismatch")
	}
}

func TestExpiredJWT(t *testing.T) {
	dur, err := time.ParseDuration("1s")
	if err != nil {
		t.Fatal(err)
	}
	testUser := uuid.New()

	token, err := MakeJWT(testUser, "Pass234", dur)
	if err != nil {
		t.Fatal(err)
	}

	sleep, err := time.ParseDuration("2s")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(sleep)

	returnedUser, err := ValidateJWT(token, "Pass234")
	if returnedUser != uuid.Nil || err == nil {
		t.Fail()
	}

}

func TestIncorrectSecretJWT(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		t.Fatal(err)
	}
	testUser := uuid.New()

	token, err := MakeJWT(testUser, "Password", dur)
	if err != nil {
		t.Fatal(err)
	}

	returnedUser, err := ValidateJWT(token, "Pa$$word")
	if returnedUser != uuid.Nil || err == nil {
		t.Fail()
	}
}
