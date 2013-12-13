package fyrirtaekjaskra

import (
	"fmt"
	"testing"
)

func TestParseISAT(t *testing.T) {
	const input = `

                                        62020 Ráðgjafarstarfsemi á sviði upplýsingatækni

                                        (Aðal)
                                        56100 Veitingastaðir

                                        (Aukanúmer)`

	expected := [2]ISATType{
		{62020, "Ráðgjafarstarfsemi á sviði upplýsingatækni", true},
		{56100, "Veitingastaðir", false},
	}

	got, err := ParseISATTypes(input)

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", expected) {
		t.Errorf("expected %v, got %v.", expected, got)
	}
}
