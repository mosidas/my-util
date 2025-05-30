package password

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// パスワードポリシー
const (
	PolicyAllChars = iota // すべての文字種を使用
	PolicyAlphaNum        // 英数字のみ
)

// 長さlengthのパスワードを生成する
// policyに応じて使用する文字種を切り替える
// upper, lower, numberをそれぞれ1文字以上含む
// PolicyAllCharsの場合はsymbolも1文字以上含む
func MakePassword(length int, policy int) (string, error) {
	if length < 4 {
		return "", errors.New("password length must be at least 4")
	}

	const (
		uppers  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowers  = "abcdefghijklmnopqrstuvwxyz"
		numbers = "0123456789"
		symbols = "!#$%&'()*+,-./:;<=>?@[]^_`{|}~"
	)

	upper, err := randomChar(uppers)
	if err != nil {
		return "", err
	}
	lower, err := randomChar(lowers)
	if err != nil {
		return "", err
	}
	number, err := randomChar(numbers)
	if err != nil {
		return "", err
	}

	var symbol rune
	var chars string
	var requiredChars int

	if policy == PolicyAllChars {
		symbol, err = randomChar(symbols)
		if err != nil {
			return "", err
		}
		chars = uppers + lowers + numbers + symbols
		requiredChars = 4
	} else if policy == PolicyAlphaNum {
		chars = uppers + lowers + numbers
		requiredChars = 3
	} else {
		return "", errors.New("invalid policy")
	}

	charsLen := len(chars)
	tmp := make([]byte, length-requiredChars)
	for i := 0; i < length-requiredChars; i++ {
		index, err := randInt(charsLen)
		if err != nil {
			return "", err
		}
		tmp[i] = chars[index]
	}

	var password string
	if policy == PolicyAllChars {
		password = string(upper) + string(lower) + string(number) + string(symbol) + string(tmp)
	} else {
		password = string(upper) + string(lower) + string(number) + string(tmp)
	}
	return shuffle(password), nil
}

func randomChar(set string) (rune, error) {
	index, err := randInt(len(set))
	if err != nil {
		return 0, err
	}
	return rune(set[index]), nil
}

// 0からn-1までのランダムな整数を返す
func randInt(n int) (int, error) {
	max := big.NewInt(int64(n))
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return int(r.Int64()), nil
}

// 文字列をランダムに並び替える
func shuffle(s string) string {
	r := []rune(s)
	for i := range r {
		j, err := randInt(i + 1)
		if err != nil {
			return s // Return the original string if an error occurs
		}
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
