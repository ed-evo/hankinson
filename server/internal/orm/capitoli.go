package orm

import (
	"fmt"
	"slices"

	"github.com/ed-evo/hankinson/server/internal/utils"
)

func provideCapitoli(base []int) []int {
	var capitoli []int
	if len(base) > 0 {
		capitoli = slices.Clone(base)
	} else {
		for key, _ := range DomandeByCapitolo {
			capitoli = append(capitoli, key)
		}
	}
	utils.ShuffleSlice(capitoli)
	return capitoli
}

func RandomDomandeIds(base []int, count int) ([]int, error) {
	var result []int
	seen := make(map[int]bool)
	capitoli := provideCapitoli(base)
	size := len(capitoli)
	for i := range max(1, count) {
		capitolo := capitoli[i%size]
		domandeIDs, ok := DomandeByCapitolo[capitolo]
		if !ok {
			return nil, fmt.Errorf("Errore capitolo %d non esiste.", capitolo)
		}
		for {
			domandaID, ok := utils.SelectRandomFromList(domandeIDs)
			if !ok {
				return nil, fmt.Errorf("Error select domanda per capitolo %d", capitolo)
			}
			if seen[domandaID] {
				continue
			}
			seen[domandaID] = true
			result = append(result, domandaID)
			break
		}
	}
	utils.ShuffleSlice(result)
	return result, nil
}
