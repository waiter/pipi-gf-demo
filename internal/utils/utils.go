package utils

import "math/rand"

func RandString(len int) string {
	bytes := make([]byte, len)
    for i := 0; i < len; i++ {
			b := rand.Intn(62)
			if b < 10 { // 0-9
				bytes[i] = byte(b + 48)
			} else if b < 36 { // A-Z
				bytes[i] = byte(b + 55)
			} else { // a-z
				bytes[i] = byte(b + 61)
			}
    }
    return string(bytes)
}