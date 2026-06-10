package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// ParseBody accepts the request body and unmarshals it into the provided interface
// the reason why interface is used is to allow any type to be passed in
// You must pass a pointer (&user)
func ParseBody(r *http.Request, x interface{}) {
	if body, err := io.ReadAll(r.Body); err == nil {
		//Take the raw JSON and turn it into a Go struct
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}


func ParseAmount(raw interface{}) (int64, error) {
	switch v := raw.(type) {
	case int64:
		if v <= 0 {
			return 0, fmt.Errorf("amount must be greater than zero")
		}
		return v, nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil || f <= 0 {
			return 0, fmt.Errorf("invalid amount: %s", v)
		}
		return int64(f), nil
	default:
		return 0, fmt.Errorf("unsupported amount type: %T", raw)
	}
}

func GenerateOTP() string {
	// Use crypto/rand for better randomness
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return fmt.Sprintf("%06d", r.Intn(1000000))
}