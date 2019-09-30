package api

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/andreylm/graphql-demo/api/errors"
)

// Screenshot - Screenshot sctruct
type Screenshot struct {
	ID      int    `json:"id"`
	VideoID int    `json:"videoId"`
	URL     string `json:"url"`
}

// User - User sctruct
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Video - Video sctruct
type Video struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	UserID      int           `json:"user"`
	URL         string        `json:"url"`
	CreatedAt   time.Time     `json:"createdAt"`
	Screenshots []*Screenshot `json:"screenshots"`
	Related     []*Video      `json:"related"`
}

// MarshalID - marshaling id
func MarshalID(id int) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", id)))
	})
}

// UnmarshalID - unmarshaling id
func UnmarshalID(v interface{}) (int, error) {
	id, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("ids myst be strings")
	}
	i, err := strconv.Atoi(id)
	return int(i), err
}

// MarshalTimestamp - marshals time
func MarshalTimestamp(t time.Time) graphql.Marshaler {
	timestamp := t.Unix() * 1000
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(timestamp, 10))
	})
}

// UnmarshalTimestamp - unmarshal timestamp
func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(int); ok {
		return time.Unix(int64(tmpStr), 0), nil
	}
	return time.Time{}, errors.TimeStampError
}
