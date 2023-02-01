package xor

func Sum[T int | int64 | int8 | int32 | byte](a T, b T) T {
	return a ^ b
}

func EncodeText(text string, key string) []byte {
	return Encode([]byte(text), key)
}

func Encode(bs []byte, key string) []byte {
	var result []byte
	for i, c := range bs {
		result = append(result, Sum(c, key[i%len(key)]))
	}

	return result
}
