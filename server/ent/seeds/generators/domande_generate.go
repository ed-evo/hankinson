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

	seedsContainer, err := seeds.GetSeeds()

	if err != nil {
		log.Fatalf("Errore lettura dati: %v", err)
	}

	capitoli := seedsContainer.Capitoli
	domande := seedsContainer.Domande
	argomennti := seedsContainer.Argomenti
	spiegazioni := seedsContainer.Spiegazioni

	slices.SortFunc(capitoli, func(a, b seeds.Capitolo) int {
		return cmp.Compare(a.ID, b.ID)
	})

	slices.SortFunc(domande, func(a, b seeds.Domanda) int {
		return cmp.Compare(a.Numero, b.Numero)
	})

	slices.SortFunc(argomennti, func(a, b seeds.Argomento) int {
		return cmp.Compare(a.ID, b.ID)
	})

	slices.SortFunc(spiegazioni, func(a, b seeds.Spiegazione) int {
		return cmp.Compare(a.Numero, b.Numero)
	})

	log.Printf("Capitoli %v, Domande %v, Argomenti %v, Spiegazioni %v", len(capitoli), len(domande), len(argomennti), len(spiegazioni))

	// 1. Definiamo le funzioni custom usate nel file SQL (.escape)
	funcMap := template.FuncMap{
		"escape": func(s string) string {
			return strings.ReplaceAll(s, "'", "''")
		},
	}

	dataDomande := struct {
		Capitoli  []seeds.Capitolo
		Argomenti []seeds.Argomento
		Domande   []seeds.Domanda
	}{
		Capitoli:  capitoli,
		Argomenti: argomennti,
		Domande:   domande,
	}

	generateSql(seeds.DomandeSqlTemplate, dataDomande, "seeds_domande.sql", funcMap)

	dataSpiegazioni := struct {
		Spiegazioni []seeds.Spiegazione
	}{
		Spiegazioni: spiegazioni,
	}
	generateSql(seeds.SpiegazioniSqlTemplate, dataSpiegazioni, "seeds_spiegazioni.sql", funcMap)

	log.Println("Done")
}

func generateSql(
	srcTemplate string,
	data any,
	dstFile string,
	funcMap template.FuncMap,
) {
	tpl, err := template.New("sql").Funcs(funcMap).Parse(srcTemplate)
	if err != nil {
		log.Fatalf("Errore creazione template: %v", err)
	}

	var queryBuffer bytes.Buffer

	if err := tpl.Execute(&queryBuffer, data); err != nil {
		log.Fatalf("errore %v", err)
	}

	if err := os.WriteFile(dstFile, queryBuffer.Bytes(), 0644); err != nil {
		log.Fatalf("errore scrittura file %v", err)
	}
}
