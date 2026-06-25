package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"
)

var generator *shortid.Shortid

var generatorSeed int64 = 1000

type clientError struct {
	ID            string `json:"id"`
	MessageToUser string `json:"messageToUser"`
	DeveloperInfo string `json:"developerInfo"`
	Error         string `json:"error"`
	StatusCode    int    `json:"statusCode"`
	IsClientError bool   `json:"isClientError"`
}

func init() {

	n, err := rand.Int(rand.Reader, big.NewInt(generatorSeed))

	if err != nil {
		logrus.Panicf("failed to initialize utils with random seed: %+v", err)
		return
	}

	g, err := shortid.New(1, shortid.DefaultABC, n.Uint64())
	if err != nil {
		logrus.Panicf("failed to initialize utils with random seed: %+v", err)
	}

	generator = g
}

func ParseBody(r io.Reader, body interface{}) error {
	err := json.NewDecoder(r).Decode(body)
	if err != nil {
		return err
	}
	return nil
}

func EncodeJSONBody(w http.ResponseWriter, body interface{}) error {
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		return err
	}
	return nil
}

func ResponseJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			logrus.Errorf("failed to response json with error: %+v", err)
		}
	}
}

func newclientError(statusCode int, messageToUser string, err error, additionalInfoDev ...string) *clientError {
	additionalInfo := strings.Join(additionalInfoDev, "/n")

	if additionalInfo == "" {
		additionalInfo = messageToUser
	}

	errorId, err := generator.Generate()

	if err != nil {
		logrus.Panicf("failed to generate errorId: %+v", err)
	}

	var errString string
	if err != nil {
		errString = err.Error()
	}

	return &clientError{
		errorId,
		messageToUser,
		additionalInfo,
		errString,
		statusCode,
		true,
	}
}

func ResponseError(w http.ResponseWriter, statusCode int, err error, messageToUser string, additionalInfoDev ...string) {
	logrus.Errorf("status: %d, message: %s, err: %+v ", statusCode, messageToUser, err)
	clientError := newclientError(statusCode, messageToUser, err, additionalInfoDev...)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientError); err != nil {
		logrus.Errorf("Failed to send error to caller with error: %+v", err)
	}
}

func HashPassword(password string) (string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil{
		return "", fmt.Errorf("failed to hash password: %w", err)
	} 
	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) error{
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashString(toHash string) string {
    sha := sha512.New()
	sha.Write([]byte(toHash))
	return hex.EncodeToString(sha.Sum(nil))		
}