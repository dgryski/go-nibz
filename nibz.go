// Package nibz compresses 4-5 byte values into 16-bits
/*
	Port of https://github.com/isometric/BucketCompressionTrick
*/
package nibz

import (
	"errors"
)

var lookupFwd = make([]uint16, 1<<16)
var lookupRev = make([]uint16, 1<<12)

func init() {
	var index uint16
	for i := 0; i < 1<<4; i++ {
		for j := 0; j <= i; j++ {
			for k := 0; k <= j; k++ {
				for l := 0; l <= k; l++ {
					sorted := uint16(i)<<(0*4) | uint16(j)<<(1*4) | uint16(k)<<(2*4) | uint16(l)<<(3*4)
					lookupFwd[sorted] = index
					lookupRev[index] = sorted
					index++
				}
			}
		}
	}
}

// 4 element sorting network
var sortingNetwork = [5][2]int{{0, 1}, {2, 3}, {0, 2}, {1, 3}, {1, 2}}

func sort(data *[4]byte) {
	for _, v := range sortingNetwork[:] {
		i, j := v[0], v[1]
		if data[i] < data[j] {
			data[i], data[j] = data[j], data[i]
		}
	}
}

var errValueTooLarge = errors.New("nibz: value too large")

func Compress(data [4]byte) (uint16, error) {

	for _, v := range data[:] {
		if v >= 1<<5 {
			return 0, errValueTooLarge
		}
	}

	sort(&data)

	sorted := uint16(data[0]>>1)<<(0*4) | uint16(data[1]>>1)<<(1*4) | uint16(data[2]>>1)<<(2*4) | uint16(data[3]>>1)<<(3*4)
	code := uint16(lookupFwd[sorted])<<4 | uint16(data[0]&1)<<0 | uint16(data[1]&1)<<1 | uint16(data[2]&1)<<2 | uint16(data[3]&1)<<3

	return code, nil
}

func Decompress(code uint16) ([4]byte, error) {

	if int(code>>4) >= len(lookupRev) {
		return [4]byte{}, errValueTooLarge
	}

	sorted := lookupRev[code>>4]

	var data [4]byte
	data[0] = byte(sorted>>(0*4)&0xf)<<1 | byte(code&(1<<0))>>0
	data[1] = byte(sorted>>(1*4)&0xf)<<1 | byte(code&(1<<1))>>1
	data[2] = byte(sorted>>(2*4)&0xf)<<1 | byte(code&(1<<2))>>2
	data[3] = byte(sorted>>(3*4)&0xf)<<1 | byte(code&(1<<3))>>3

	return data, nil
}
