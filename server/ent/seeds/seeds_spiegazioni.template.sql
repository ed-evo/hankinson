BEGIN TRANSACTION;
-- INSERIMENTO SPIEGAZIONI
{{range .Spiegazioni}}
insert into spiegazioni (numero_domanda, spiegazione, focus_linguistico, regola_chiave)
VALUES ({{.Numero}}, '{{escape .Spiegazione}}', '{{escape .FocusLinguistico}}', '{{escape .RegolaChiave}}')
ON CONFLICT(numero_domanda) DO UPDATE SET
    numero_domanda = excluded.numero_domanda,
    spiegazione = excluded.spiegazione,
    focus_linguistico = excluded.focus_linguistico,
    regola_chiave = excluded.regola_chiave;
{{end}}

COMMIT;