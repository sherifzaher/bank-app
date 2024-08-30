package api

import "net/http"

const (
	NotFound = "sql: no rows in result set"
)

func GetError(err error) (string, int) {
	switch err.Error() {
	case NotFound:
		return "Item not found", http.StatusNotFound
	}
	return "Unexpected error", http.StatusInternalServerError
}
