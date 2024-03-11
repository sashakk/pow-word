package quote

import "math/rand"

var quotes = []string{
	"Life can only be understood backwards; but it must be lived forwards.",
	"People demand freedom of speech as a compensation for the freedom of thought which they seldom use.",
	"Life is not a problem to be solved, but a reality to be experienced.",
}

func GetQuote() string {
	return quotes[rand.Intn(len(quotes))]
}
