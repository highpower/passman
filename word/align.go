package word

func alignedSize(size int) int {
	x := 8
	for i := 0; x < size; i++ {
		x <<= 1
	}
	return x
}
