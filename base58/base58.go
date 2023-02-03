package base58

import (
	"bytes"
	"math/big"
)

var base58 = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Encode(str string) string {
	sb := []byte(str)
	sbt := big.NewInt(0).SetBytes(sb)

	var modSlice []byte
	for sbt.Cmp(big.NewInt(0)) > 0 {
		mod := big.NewInt(0)
		sbt.DivMod(sbt, big.NewInt(58), mod)
		modSlice = append(modSlice, base58[mod.Int64()])
	}

	for _, elem := range sb {
		if elem != 0 {
			break
		} else if elem == 0 {
			modSlice = append(modSlice, byte('1'))
		}
	}

	return string(reverseBytes(modSlice))
}

func reverseBytes(bytes []byte) []byte {
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i]
	}
	return bytes
}

func Decode(str string) []byte {
	strByte := []byte(str)
	ret := big.NewInt(0)
	for _, byteElem := range strByte {
		index := bytes.IndexByte(base58, byteElem)
		ret.Mul(ret, big.NewInt(58))
		ret.Add(ret, big.NewInt(int64(index)))
	}

	return ret.Bytes()
}
