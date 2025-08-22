package sat

func abs(x int) int {
	switch {
	case x > 0:
		return x
	case x < 0:
		return -x
	default:
		return 0
	}
}

func sign(x int) int {
	switch {
	case x > 0:
		return +1
	case x < 0:
		return -1
	default:
		return 0
	}
}
