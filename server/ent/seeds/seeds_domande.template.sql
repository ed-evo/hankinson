BEGIN TRANSACTION;
-- INSERIMENTO ARGOMENTI (Se l'ID esiste già, aggiorna il nome)
{{range .Argomenti}}
INSERT INTO argomenti (id, nome) 
VALUES ({{.ID}}, '{{escape .Nome}}') 
ON CONFLICT(id) DO UPDATE SET nome = excluded.nome;
{{end}}
-- INSERIMENTO DOMANDE (Se il numero esiste già, aggiorna tutti i campi)
{{range .Domande}}
INSERT INTO domande (numero, testo, is_true, immagine, pagina_quiz, id_blocco) 
VALUES ({{.Numero}}, '{{escape .Testo}}', {{.IsTrue}}, {{if .Immagine}}'{{escape .Immagine}}'{{else}}NULL{{end}}, {{.PaginaQuiz}}, {{.IDBlocco}}) 
ON CONFLICT(numero) DO UPDATE SET 
    testo = excluded.testo, 
    is_true = excluded.is_true, 
    immagine = excluded.immagine, 
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