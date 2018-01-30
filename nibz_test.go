package nibz

import "testing"

func TestSort(t *testing.T) {

	const max = 1 << 5

	for i := 0; i < max; i++ {
		for j := 0; j < max; j++ {
			for k := 0; k < max; k++ {
				for l := 0; l < max; l++ {
					var data = [4]byte{byte(i), byte(j), byte(k), byte(l)}

					sort(&data)

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
