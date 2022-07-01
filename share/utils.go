package share

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

func ExitFatal(format string, args ...interface{}) {
	errorMessage := fmt.Sprintf(format, args...)
	log.Fatalf("fatal error: %s", errorMessage)
}

func MustGetEnv(key string) string {
	var temp = os.Getenv(key)
	if temp == "" {
		ExitFatal("no value found in environment for key: %s", key)
	}
	return temp
}

func GetAPIVersionFromHeader(r *http.Request) (string, error) {
	headers := r.Header
	version := headers.Get("Accept-version")
	if version == "" {
		return "", NewApiErr(http.StatusBadRequest, "Accept-version header missing from request", "")
	}
	return version, nil
}

// FileWithLineNum returns the file name and line number of the function that it was called by.
// It can be used to determine the approximate location of an error
func FileWithLineNum() string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	} else {
		return "could not find filename and line num"
	}
}
