package utils

import (
	"cmp"
	"iter"
	"maps"
	"math/rand/v2"
	"slices"
)

func Clamp[T cmp.Ordered](val, minBound, maxBound T) T {
	return max(minBound, min(val, maxBound))
}

func SelectRandomFromList[T any](items []T) (T, bool) {
	var zero T
	if len(items) == 0 {
		return zero, false
	}
	return items[rand.N(len(items))], true
}

func SelectRandomFromSeq[T any](s iter.Seq[T]) (T, bool) {
	return SelectRandomFromList(slices.Collect(s))
}

func SelectRandoMapKey[K comparable, V any](m map[K]V) (K, bool) {
	return SelectRandomFromSeq(maps.Keys(m))
}
