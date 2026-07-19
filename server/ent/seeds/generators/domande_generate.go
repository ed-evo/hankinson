package main

import (
	"bytes"
	"cmp"
	"log"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/ed-evo/hankinson/server/ent/seeds"
)

func main() {

	goFile := os.Getenv("GOFILE")
	if goFile == "" {
		log.Fatal("this generator must be run via 'go generate'")
	}

	capitoli, domande, argomennti, err := seeds.GetSeeds()

	slices.SortFunc(capitoli, func(a, b seeds.Capitolo) int {
		return cmp.Compare(a.ID, b.ID)
	})

	slices.SortFunc(domande, func(a, b seeds.Domanda) int {
		return cmp.Compare(a.Numero, b.Numero)
	})

	slices.SortFunc(argomennti, func(a, b seeds.Argomento) int {
		return cmp.Compare(a.ID, b.ID)
	})

	if err != nil {
		log.Fatalf("Errore lettura dati: %v", err)
	}

	log.Printf("Capitoli %v, Domande %v, Argomenti %v", len(capitoli), len(domande), len(argomennti))

	// 1. Definiamo le funzioni custom usate nel file SQL (.escape)
	funcMap := template.FuncMap{
		"escape": func(s string) string {
			return strings.ReplaceAll(s, "'", "''")
		},
	}

	tpl, err := template.New("sql").Funcs(funcMap).Parse(seeds.DomandeSqlTemplate)
	if err != nil {
		log.Fatalf("Errore creazione template: %v", err)
	}

	var queryBuffer bytes.Buffer
	data := struct {
		Capitoli  []seeds.Capitolo
		Argomenti []seeds.Argomento
		Domande   []seeds.Domanda
	}{
		Capitoli:  capitoli,
		Argomenti: argomennti,
		Domande:   domande,
	}
	if err := tpl.Execute(&queryBuffer, data); err != nil {
		log.Fatalf("errore %v", err)
	}

	if err := os.WriteFile("seeds_domande.sql", queryBuffer.Bytes(), 0644); err != nil {
		log.Fatalf("errore scrittura file %v", err)
	}
	log.Println("Done")
}
