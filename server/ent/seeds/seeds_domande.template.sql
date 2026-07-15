BEGIN TRANSACTION;
-- INSERIMENTO CAPITOLI (Se l'id esiste già, aggiorna tutti i campi)
{{range .Capitoli}}
INSERT INTO capitoli (id, nome, min_numero_domanda, max_numero_domanda, totale_domande) 
VALUES ({{.ID}}, '{{escape .Nome}}', {{.MinNumeroDomanda}}, {{.MaxNumeroDomanda}}, {{.TotaleDomande}}) 
ON CONFLICT(id) DO UPDATE SET 
    nome = excluded.nome,
    min_numero_domanda = excluded.min_numero_domanda,
    max_numero_domanda = excluded.max_numero_domanda,
    totale_domande = excluded.totale_domande;
{{end}}
-- INSERIMENTO ARGOMENTI (Se l'ID esiste già, aggiorna il nome)
{{range .Argomenti}}
INSERT INTO argomenti (id, nome) 
VALUES ({{.ID}}, '{{escape .Nome}}') 
ON CONFLICT(id) DO UPDATE SET nome = excluded.nome;
{{end}}
-- INSERIMENTO DOMANDE (Se il numero esiste già, aggiorna tutti i campi)
{{range .Domande}}
INSERT INTO domande (numero, testo, is_true, immagine, id_capitolo, pagina_quiz, id_blocco) 
VALUES ({{.Numero}}, '{{escape .Testo}}', {{.IsTrue}}, {{if .Immagine}}'{{escape .Immagine}}'{{else}}NULL{{end}}, {{.IDCapitolo}}, {{.PaginaQuiz}}, {{.IDBlocco}}) 
ON CONFLICT(numero) DO UPDATE SET 
    testo = excluded.testo, 
    is_true = excluded.is_true, 
    immagine = excluded.immagine, 
    id_capitolo = excluded.id_capitolo,
    pagina_quiz = excluded.pagina_quiz, 
    id_blocco = excluded.id_blocco;
{{end}}
-- ASSOCIAZIONI ARGOMENTI_DOMANDE (Se la coppia esiste già, ignora per evitare duplicati)
{{range $domanda := .Domande}}{{range $argId := $domanda.Argomenti}}
INSERT INTO argomenti_domande (argomento_id, domanda_id)
VALUES ({{$argId}}, {{$domanda.Numero}}) 
ON CONFLICT(argomento_id, domanda_id) DO NOTHING;
{{end}}{{end}}

COMMIT;