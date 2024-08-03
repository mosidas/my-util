package password

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// 長さlengthのパスワードを生成する
// upper, lower, number, symbolをそれぞれ1文字以上含む
func MakePassword(length int) (string, error) {
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
	symbol, err := randomChar(symbols)
	if err != nil {
		return "", err
	}

	chars := uppers + lowers + numbers + symbols
	charsLen := len(chars)
	tmp := make([]byte, length-4)
	for i := 0; i < length-4; i++ {
		index, err := randInt(charsLen)
		if err != nil {
			return "", err
		}
		tmp[i] = chars[index]
	}

	password := string(upper) + string(lower) + string(number) + string(symbol) + string(tmp)
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
