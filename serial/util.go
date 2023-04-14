package serial

import "time"

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

func runWithTimeout(timeout time.Duration, f func() interface{}, errorMessage ...string) interface{} {
	done := make(chan interface{})
	var data interface{} = nil

	go func() {
		done <- f()
	}()

	select {
	case data = <-done:
		return data
	case <-time.After(timeout):
		if len(errorMessage) > 0 {
			panic(errorMessage[0])
		}

		panic("Timeout reached")
	}
}
