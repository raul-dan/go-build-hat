package serial

func checksum(buff []byte) uint32 {
	var sum uint32 = 1

	for i := 0; i < len(buff); i++ {
		if (sum & 0x80000000) != 0 {
			sum = (sum << 1) ^ 0x1D872B41
		} else {
			sum = sum << 1
		}

		sum = (sum ^ uint32(buff[i])) & 0xFFFFFFFF
	}

	return sum
}
