package errors

import (
	"errors"
	"fmt"
	"log"
	"runtime"

	"github.com/lib/pq"
)

var (
	// ServerError - server error
	ServerError = GenerateError("Something went wrong! Please try again later")
	// UserNotExist - user not exist error
	UserNotExist = GenerateError("User not Exist")
	// UnauthorisedError - unauthorised error
	UnauthorisedError = GenerateError("You are not authorised to perform this action")
	// TimeStampError - invalid timestamp
	TimeStampError = GenerateError("time should be a unix timestamp")
	// InternalServerError - internal server error
	InternalServerError = GenerateError("internal server error")
)

// GenerateError - generates error
func GenerateError(err string) error {
	return errors.New(err)
}

// IsForeignKeyError - checks if error is foreign key error
func IsForeignKeyError(err error) bool {
	pgErr := err.(*pq.Error)
	return pgErr.Code == "23503"
}

// DebugPrintf - debug error
func DebugPrintf(err error, args ...interface{}) string {
	programCounter, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(programCounter)
	msg := fmt.Sprintf("[%s: %s %d] %s, %s", file, fn.Name(), line, err, args)
	log.Println(msg)
	return msg
}
