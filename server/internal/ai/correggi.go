package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/attivitaquesitoesame"
	"github.com/ed-evo/hankinson/server/ent/esame"
	"google.golang.org/genai"
)

type EsameCorrezione struct {
	Model string
	Promosso bool
	Testo string
	Metadata string
}

const correzioneSystemPrompt = `## RUOLO
Sei un istruttore di scuola guida esperto, chiaro, empatico e incoraggiante. 

## OBIETTIVO
Analizzare il report/log di un esame di teoria della Patente B sostenuto da uno studente e generare un feedback didattico completo, chiaro e personalizzato.

---

## ⛔ DIRETTIVE TASSATIVE DI FORMATTAZIONE (PREAMBLE ELIMINATION)
1. **NESSUN PREAMBOLO O SALUTO:** Inizia la risposta **DIRETTAMENTE** con il primo intestazione '### 1. Esito Generale & Gestione del Tempo'.
2. **NESSUN TESTO DI INTRODUZIONE O CONCLUSIONE:** Non aggiungere mai frasi come *"Ecco l'analisi dell'esame:"*, *"Certamente!"*, o *"Buona fortuna per la prossima volta!"* all'inizio o alla fine del testo.
3. **OUTPUT RIGIDAMENTE IN MARKDOWN:** Rispondi esclusivamente organizzando il testo nelle 3 sezioni indicate sotto.

---

## STRUTTURA DELL'INPUT
Riceverai i dati dell'esame sotto forma di tabella riassuntiva e/o registro attività (log cronologico dei singoli quesiti), contenenti:
- Testo della domanda del quiz
- Risposta data dall'alunno vs Risposta corretta (V/F)
- Tempo di risposta in millisecondi (ms)
- Eventuali cambi di idea o rettifiche durante lo svolgimento

---

## REGOLE E DIRETTIVE DI ANALISI

1. **Soglie di Valutazione Tempi:**
   - **Risposta D'Impulso (< 3000 ms / 3s):** Se lo studente risponde in meno di 3 secondi su domande di comprensione o incroci, evidenzia il rischio di disattenzione.
   - **Forte Incertezza (> 15000 ms / 15s):** Se lo studente impiega più di 15 secondi, segnala il quesito come punto critico/da ripassare, anche se la risposta finale è corretta.

2. **Criteri di Valutazione Errori:**
   - Spiega in modo chiaro e semplice PERCHÉ la risposta corretta è quella.
   - Spiega il trabocchetto o l'errore concettuale commesso dallo studente.
   - Fai riferimento diretto alle norme del Codice della Strada in modo accessibile (senza citare articoli di legge complessi).

3. **Criteri di Riformattazione Tempi:**
   - Converti sempre i millisecondi (ms) in secondi (s) per rendere il report leggibile all'utente (es. 12202ms -> ~12 secondi).

---

## FORMATO DELLO OUTPUT
Il primo carattere della tua risposta deve essere il simbolo '#'. Rispondi organizzando il testo **tassativamente** nelle seguenti 3 sezioni Markdown:

### 1. Esito Generale & Gestione del Tempo
- **Esito finale:** [PROMOSSO / BOCCIATO] con [N] errori su 10 (Massimo errori ammessi: 3).
- **Analisi del Ritmo di Esame:** Commento sintetico sulla gestione del tempo complessivo e sul livello di sicurezza mostrato.

### 2. Analisi Dettagliata dei Quesiti
- **Errori Commessi:**
  * **Quesito N. [ID]** - *"[Testo della domanda]"*
    - **La tua risposta:** [Vero/Falso] | **Risposta corretta:** [Vero/Falso] | **Tempo:** [X]s
    - **Spiegazione didattica:** [Perché la risposta corretta è quella e dove si trova l'errore/trabocchetto]
- **Punti d'Incertezza (Risposte corrette ma lente):**
  * **Quesito N. [ID]** - Tempo impiegato: [X]s. [Breve nota sul perché questo argomento richiede un ripasso]

### 3. Consigli Pratici & Prossimi Passi
- **Macro-area critica:** [Indica l'argomento dove si sono concentrati gli dubbi o gli errori]
- **3 Azioni consigliate:**
  1. [Consiglio pratico 1]
  2. [Consiglio pratico 2]
  3. [Consiglio pratico 3]

---

## TONO DI VOCE
- Professionale, motivante e costruttivo.
- Accessibile e mai burocratico.
- Se lo studente è bocciato, mantieni un approccio incoraggiante focalizzato sul miglioramento.
`

var esameSummaryTemplate = template.Must(template.New("summary-table").Parse(`### Tabella riassuntiva dell'esame sostenuto
{{ with .Esame }}
Esame ID: {{ .ID }}
Numero Quesiti: {{ .NumeroQuesiti }}
Errori ammessi: {{ .MaxErrori }}
Tempo massimo: {{ .MinutiDisponibili }} minuti
{{ end }}
| Quesito ID | Domanda | Risposta Corretta | Risposta Alunno | Esito | Tempo Impiegato |
| ----- | ----- | ----- | ----- | ----- | ----- |
{{- range .Summaries }}
| {{ .ID }} | {{ .Domanda }} | {{ .RispostaCorretta }} | {{ .Risposta }} | {{ .Esito }} | {{ .TempoImpiegatoMs }}ms |
{{- end }}
`))

type EsameSummaryInput struct {
	Esame *ent.Esame
	Summaries map[int]*QuesitoResponse
}
func compileEsameSummary (input *EsameSummaryInput) (string, error) {
	var buf bytes.Buffer
	if err := esameSummaryTemplate.Execute(&buf, input); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type QuesitoResponse struct {
	ID int
	Domanda string
	RispostaCorretta string
	Risposta string
	Esito string
	TempoImpiegatoMs int

}

func vfy (vbool *bool) string {
	if (vbool == nil) {
		return "N/D"
	} else if (*vbool == true) {
		return "V"
	} else {
		return "F"
	}
}

func (ai *Gem) Correggi(ctx context.Context, db *ent.Client, id int) (*EsameCorrezione, error) {

	e, err := db.Esame.Get(ctx, id)
	if err != nil {
		log.Printf("Erore lettura esame %d", id)
		return nil, err
	}

	log.Print(e)

	listaAttivita, err := db.Esame.Query().
	Where(esame.ID(id)).
	QueryQuesiti().QueryAttivita().
	WithQuesitoEsame(func(qeq *ent.QuesitoEsameQuery) {
		qeq.WithDomandaOriginale()
	}).All(ctx)

	if err != nil {
		log.Printf("Errore lettura attivita per esame %d", id)
		return nil, err
	}

	var currentQuesito *ent.QuesitoEsame
	tempoQuesitoCorrente := 0
	countSbagliate := 0

	var risultati = make(map[int]*QuesitoResponse)

	
	var parts []*genai.Part

	parts = append(parts, genai.NewPartFromText("### Registro attività dell'esame sostenuto dall'alunno"))

	for _, a := range listaAttivita {
		q := a.Edges.QuesitoEsame
		d := q.Edges.DomandaOriginale
		if currentQuesito != q {
			if qr, ok := risultati[q.ID]; !ok {

				var esito string
				if q.RispostaFinale != nil && *q.RispostaFinale == d.IsTrue {
					esito = "CORRETTA"
				} else {
					esito  = "SBAGLIATA"
					countSbagliate += 1
				}
				risultati[q.ID] = &QuesitoResponse{
					ID: q.ID,
					Domanda: d.Testo,
					RispostaCorretta: vfy(&d.IsTrue),
					Risposta: vfy(q.RispostaFinale),
					Esito: esito,
					TempoImpiegatoMs: tempoQuesitoCorrente,
				}
			} else {
				qr.TempoImpiegatoMs += tempoQuesitoCorrente
			}
			currentQuesito = q
			tempoQuesitoCorrente = 0
			if d.Immagine != nil {
				img, _ := toImagePart(d)
				parts = append(parts, img)
			}
			parts = append(parts, genai.NewPartFromText(fmt.Sprintf("Quesito N. %d\n[DOMANDA] %s\n", q.ID, d.Testo)))
		}
		tempoQuesitoCorrente = tempoQuesitoCorrente + a.DurataMs
		var sb strings.Builder
		switch a.Tipo {
		case attivitaquesitoesame.TipoRisposta:
			var resUtente string
			if a.RispostaData == nil {
				resUtente = "blank"
			} else if *a.RispostaData == true {
				resUtente = "Vero"
			} else {
				resUtente = "Falso"
			}
			sb.WriteString("[RISPOSTA] ")
			sb.WriteString(resUtente)
		case attivitaquesitoesame.TipoSalta:
			sb.WriteString("[AZIONE] Salta")
		case attivitaquesitoesame.TipoProssimo:
			sb.WriteString("[AZIONE] prossima")
		}
		sb.WriteString(" (")
		sb.WriteString(strconv.Itoa(a.DurataMs))
		sb.WriteString("ms)\n")
		parts = append(parts, genai.NewPartFromText(sb.String()))
	}
	parts = append(parts, genai.NewPartFromText("[AZIONE] Fine esame\n"))

	if sum, err := compileEsameSummary(&EsameSummaryInput{Esame: e, Summaries: risultati}); err != nil {
		return nil, err
	} else {
		parts = append(parts, genai.NewPartFromText(sum))
	}

	userInputs := genai.NewContentFromParts(parts, genai.RoleUser)

	aiConfig := &genai.GenerateContentConfig{
		SystemInstruction: genai.Text(correzioneSystemPrompt)[0],
	}

	resp, err := ai.client.Models.GenerateContent(
		ctx, ai.model,
		[]*genai.Content{userInputs},
		aiConfig,
	)

	if err != nil {
		return nil, fmt.Errorf("gemini generation error: %w", err)
	}

	correzione := &EsameCorrezione{
		Model: ai.model,
		Promosso: countSbagliate <= e.MaxErrori,
		Testo: resp.Text(),
	}

	resp.Candidates = nil

	if bytes, err := json.MarshalIndent(resp, "", " "); err == nil {
		correzione.Metadata = string(bytes)
		
	}

	return correzione, nil
}