package quote

import "testing"

func TestGetQuote(t *testing.T) {
	t.Run("quote test", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			q := GetQuote()
			if len(q) == 0 {
				t.Errorf("empty quote")
			}
		}
	})
}
