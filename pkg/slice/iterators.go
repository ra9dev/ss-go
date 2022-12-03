package slice

func Has[T comparable](in []T, val T) bool {
	for _, v := range in {
		if v == val {
			return true
		}
	}

	return false
}
