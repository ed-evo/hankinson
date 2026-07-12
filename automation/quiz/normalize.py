
from typing import Optional
import re

import context

from custo_types import (
    Metadata,
    RispostaCorretta,
    Domanda
)


pattern = r"Quesito n°\s*(?P<quesito>\d+)\s*-\s*(?P<argomenti>.+)"
def parse_caption(caption: str) -> Optional[Metadata]:
    logger = context.getLogger()
    if not caption:
        logger.error("Tentativo di estrazione quesito e argomenti da intestazione vuota")
        raise AssertionError("Intestazione vuota")
    matches = re.search(pattern, caption)

    if not matches:
        logger.error(f"Formato errato per intestazione: {caption}")
        raise AssertionError("Errore formato intestazione")
    
    quesito = int(matches.group("quesito"))
    argomenti = [stripped for argomento in matches.group("argomenti").split(";") if (stripped := argomento.strip())]

    logger.debug(f"Quesito {quesito}: argomenti {argomenti}")
    return {
        "id_blocco": quesito,
        "argomenti": argomenti
    }

def extract_cell_text(cell, strict: bool = True) -> str:
    logger = context.getLogger()
    text = cell.get("text")
    images = cell.get("images")
    logger.debug(f"Contenuto: {text}")
    if strict and images:
        logger.error("Trovata presenza di immagini su cella di test")
        raise AssertionError("Immagini in Cella testo")

    return text

def extract_cell_image(cell, strict: bool = True) -> Optional[str]:
    logger = context.getLogger()
    text = cell.get("text")
    images = cell.get("images")
    logger.debug(f"Immagini: {images}")
    if strict:
        if text:
            logger.error("Trovato presenza di testo in cella immagine")
            raise AssertionError("Testo in cella immagine")
        if len(images) > 1:
            logger.error("Una cella immagine puo' avere al massimo una sola immagine")
            raise AssertionError("Troppe immagini per cella immagine")

    return images[0] if images else None
def parse_riga(riga, metadata: Metadata) -> Domanda:
    logger = context.getLogger()
    if not riga:
        logger.error("Riga assente")
        raise AssertionError("Riga vuota")
    
    if not riga.get("page"):
        logger.error("Pagina obbligatoria per la riga.")
        raise AssertionError("Riga senza pagina")
    
    if not riga.get("contents"):
        logger.error("Contenuti obbligatorio per le righe")
        raise AssertionError("Riga senza contenuti")
    
    cols = riga.get("contents")
    if (cols_len := len(cols)) != 4:
        logger.error(f"Trovate {cols_len} invece di 4")
        raise AssertionError("Riga non ha 4 colonne")
    
    numero_domanda = extract_cell_text(cols[0])
    testo_domanda = extract_cell_text(cols[1])
    risposta_corretta = extract_cell_text(cols[2])
    immagine = extract_cell_image(cols[3])

    return {
        "id_blocco": metadata["id_blocco"],
        "argomenti": metadata["argomenti"],
        "numero": int(numero_domanda),
        "testo": testo_domanda,
        "risposta_corretta": RispostaCorretta(risposta_corretta),
        "immagine": immagine,
        "pagina_quiz": int(riga.get("page"))
    }

def extract_quesiti(tables):
    logger = context.getLogger()
    logger.info("Inizio normalizzazione quesiti")

    domande: list[Domanda] = []

    for table in tables:
        metadata = parse_caption(table['caption'])
        logger.debug(metadata)
        # first raw is the table header
        for riga in table["rows"][1:]:
            domanda = parse_riga(riga, metadata)
            logger.info(domanda)
            domande.append(domanda)
    
    logger.info(f"Estratte {len(domande)} domande")
    return domande