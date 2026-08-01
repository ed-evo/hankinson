package utils

import (
	"iter"
	"maps"
	"math/rand/v2"
	"slices"
)

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

func ShuffleSlice[T any](s []T) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}
