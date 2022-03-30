package random

import (
	"math/rand"
	"net/http"
)

var (
	httpMethods = [...]string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
	}
	httpStatuses = [...]int{
		http.StatusContinue,
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNoContent,
		http.StatusPartialContent,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusNotModified,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusGone,
		http.StatusRequestEntityTooLarge,
		http.StatusTeapot,
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
	}
	httpVersions = [...]string{
		"HTTP/1.0",
		"HTTP/1.1",
		"HTTP/2",
	}
	userAgents = [...]string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:98.0) Gecko/20100101 Firefox/98.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12.3; rv:98.0) Gecko/20100101 Firefox/98.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:98.0) Gecko/20100101 Firefox/98.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/98.0 Mobile/15E148 Safari/605.1.15",
		"Mozilla/5.0 (Android 12; Mobile; rv:68.0) Gecko/68.0 Firefox/98.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/99.0.4844.59 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 15_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Mobile/15E148 Safari/604.1",
	}
)

// HTTPMethod returns a random HTTP method.
func HTTPMethod() string {
	return httpMethods[rand.Intn(len(httpMethods))]
}

// HTTPStatus returns a random HTTP status code.
func HTTPStatus() int {
	return httpStatuses[rand.Intn(len(httpStatuses))]
}

// HTTPVersion returns a random HTTP version.
func HTTPVersion() string {
	return httpVersions[rand.Intn(len(httpVersions))]
}

// UserAgent returns a random user agent string.
func UserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}
