package nibz

import (
	"encoding/binary"
	"testing"
)

func TestSort(t *testing.T) {

	const max = 1 << 5

	for i := 0; i < max; i++ {
		for j := 0; j < max; j++ {
			for k := 0; k < max; k++ {
				for l := 0; l < max; l++ {

					a, b, c, d := sort(byte(i), byte(j), byte(k), byte(l))

					data := [4]byte{a, b, c, d}

					var prev byte = 255
					for _, v := range data[:] {
						if v > prev {
							t.Fatalf("failed to sort: [%v,%v,%v,%v], got %v", i, j, k, l, data)
						}
						prev = v
					}
				}
			}
		}
	}
}

func TestRoundtrip(t *testing.T) {

	const max = 1 << 5

	for i := 0; i < max; i++ {
		for j := 0; j < i; j++ {
			for k := 0; k < j; k++ {
				for l := 0; l < k; l++ {
					var data = [4]byte{byte(i), byte(j), byte(k), byte(l)}

					code, _ := Compress(data)
					got, _ := Decompress(code)

					if got != data {
						t.Fatalf("Roundtrip(%v) failed: code=%v, got %v", data, code, got)
					}
				}
			}
		}
	}
}

var sink byte

func BenchmarkCompress(b *testing.B) {

	var data [4]byte

	var x uint64 = 1

	for i := 0; i < b.N; i++ {
		x = xorshiftMult64(x)
		binary.LittleEndian.PutUint32(data[:], uint32(x))
		data[0] = byte(x) & 0x1f
		data[1] = byte(x>>5) & 0x1f
		data[2] = byte(x>>10) & 0x1f
		data[3] = byte(x>>15) & 0x1f
		c, _ := Compress(data)
		sink += byte(c)
	}
}

func BenchmarkDecompress(b *testing.B) {

	var data [4]byte

	var x uint64 = 1

	for i := 0; i < b.N; i++ {
		x = xorshiftMult64(x)
		binary.LittleEndian.PutUint32(data[:], uint32(x))
		data, _ := Decompress(uint16(x))
		sink += data[0]
	}
}

func BenchmarkSort(b *testing.B) {

	var data [4]byte

	var x uint64 = 1

	for i := 0; i < b.N; i++ {
		x = xorshiftMult64(x)
		binary.LittleEndian.PutUint32(data[:], uint32(x))
		a, _, _, _ := sort(data[0], data[1], data[2], data[3])
		sink += a
	}
}

// 64-bit xorshift multiply rng from http://vigna.di.unimi.it/ftp/papers/xorshift.pdf
func xorshiftMult64(x uint64) uint64 {
	x ^= x >> 12 // a
	x ^= x << 25 // b
	x ^= x >> 27 // c
	return x * 2685821657736338717
}
