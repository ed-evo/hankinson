package utils

import (
	"cmp"
)

func Clamp[T cmp.Ordered](val, minBound, maxBound T) T {
	return max(minBound, min(val, maxBound))
}
