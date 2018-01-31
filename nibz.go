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

func sort(a, b, c, d byte) (byte, byte, byte, byte) {
	if a < b {
		a, b = b, a
	}
	if c < d {
		c, d = d, c
	}
	if a < c {
		a, c = c, a
	}
	if b < d {
		b, d = d, b
	}
	if b < c {
		b, c = c, b
	}

	return a, b, c, d
}

var errValueTooLarge = errors.New("nibz: value too large")

func Compress(data [4]byte) (uint16, error) {

	for _, v := range data[:] {
		if v >= 1<<5 {
			return 0, errValueTooLarge
		}
	}

	a, b, c, d := sort(data[0], data[1], data[2], data[3])

	sorted := uint16(a>>1)<<(0*4) | uint16(b>>1)<<(1*4) | uint16(c>>1)<<(2*4) | uint16(d>>1)<<(3*4)
	code := uint16(lookupFwd[sorted])<<4 | uint16(a&1)<<0 | uint16(b&1)<<1 | uint16(c&1)<<2 | uint16(d&1)<<3

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
