package scanner

import (
	"math/big"
	"testing"
)

func TestWeiToEth(t *testing.T) {
	tests := []struct {
		input    *big.Int
		expected string
	}{
		{big.NewInt(1e18), "1.000000"},
		{big.NewInt(5e17), "0.500000"},
		{big.NewInt(0), "0.000000"},
	}

	for _, test := range tests {
		result := weiToEth(test.input).Text('f', 6)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

