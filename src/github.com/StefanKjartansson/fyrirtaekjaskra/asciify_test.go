package fyrirtaekjaskra

import (
	"testing"
)

func TestAsciify(t *testing.T) {

	const expected = "thaedyiouoae-thaedyiouoae"
	const input = "þæðýíóúöáé-ÞÆÐÝÍÓÚÖÁÉ"

	a := Asciify(input)

	if a != expected {
		t.Errorf("Asciify: %v, expected: %v", a, expected)
	}

}
