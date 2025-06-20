package util

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const uperAlphabet = "ABCDEFGHIJKMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*max
}

func RandomPgNumeric(min, max int64, exp int32) pgtype.Numeric {
	randomInt := RandomInt(min, max)
	randomFraction := rand.Int63n(100)
	numericValue := big.NewInt(randomInt*1000 + randomFraction)
	var numeric pgtype.Numeric
	numeric.Int = numericValue
	// Set the exponent to position the decimal point (6 means 6 decimal places)
	numeric.Exp = -exp
	numeric.Valid = true
	return numeric
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUpperString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "VND"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEntryType() string {
	types := []string{"IT", "TM"}
	n := len(types)
	return types[rand.Intn(n)]
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
