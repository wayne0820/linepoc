package game

import (
	"fmt"
	"strings"
)

func GetBSResult(hash string) (bool, error) {
	hash = replaceEngChar(hash)
	ln := len(hash)
	r := hash[ln-1 : ln]

	switch r {
	case "0", "1", "2", "3", "4":
		return false, nil
	case "5", "6", "7", "8", "9":
		return true, nil
	default:
		return false, fmt.Errorf("GetBSResult parsing number error: " + r)
	}
}

func GetLuckyResult(hash string) (bool, error) {
	var b1, b2 bool
	ln := len(hash)

	r1 := hash[ln-1 : ln]
	switch r1 {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		b1 = false
	default:
		b1 = true
	}

	r2 := hash[ln-2 : ln-1]
	switch r2 {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		b2 = false
	default:
		b2 = true
	}

	if b1 != b2 {
		return true, nil
	} else {
		return false, nil
	}
}

func replaceEngChar(hash string) string {
	replacer := strings.NewReplacer("a", "", "b", "", "c", "", "d", "", "e", "", "f", "")
	return replacer.Replace(hash)
}
