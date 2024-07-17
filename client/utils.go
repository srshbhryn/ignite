package client

func hashCode(bytes []byte) int32 {
	if len(bytes) == 0 {
		return 1
	}

	h := uint32(0)
	for i := 0; i < len(bytes); i++ {
		h = 31*h + uint32(bytes[i])
	}
	return int32(h)
}
