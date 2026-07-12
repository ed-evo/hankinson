import json
from operator import itemgetter
from pathlib import Path
import context

def extract_cell_data(page_data, cell_rect):
    logger = context.getLogger()
    logger.debug(f"Estrazioni dati cella {cell_rect}")
    results = [
        item
            for item in page_data
            if cell_rect["x0"] <= item["x0"] <= cell_rect["x1"]
                and cell_rect["top"] <= item["top"] <= cell_rect["bottom"]
    ]

    logger.debug(f"Estratti {len(results)}")

    return results

def join_chars(chars):
    logger = context.getLogger()
    if not chars:
        logger.debug("chars vuoto.")
        return ""

    words = []
    current_word = []
    previous_char = None
    for char in chars:
        if char["text"] == " ":
            if current_word:
                words.append("".join(current_word))
                current_word = []
            previous_char = None
            continue
        if previous_char is not None:
            if char["top"] != previous_char["top"]:
                words.append("".join(current_word))
                current_word = []
            elif char["x0"] > previous_char["x1"]:
                words.append("".join(current_word))
                current_word = []
        current_word.append(char["text"])
        previous_char = char
    if current_word:
        words.append("".join(current_word))

    result = " ".join(words)
    logger.debug(result)
    return result

def extract_quizes(pages):
    logger = context.getLogger()
    logger.info(f"Estrazione tabelle quiz su {len(pages)} pagine")
    first_page = pages[0]
    header_rect = {
        "x0": 0,
        "x1": 600,
        "top": 0,
        "bottom": 100
    }
    logger.debug(f"aggiunto {header_rect} per incapsulare l'intestazione nella prima pagina")
    first_page["rects"].insert(0, header_rect)

    by_top_x0 = itemgetter("top", "x0")

    tables = []
    current_table = None
    bottom_text = None
    for page in pages:
        page_number = page["pageNumber"]
        chars = sorted(page["chars"], key=by_top_x0)
        rects = sorted(page["rects"], key=by_top_x0)
        images = sorted(page["images"], key=by_top_x0)
        logger.info(f"📄 Pagina: {page_number}")
        if current_table is None:
            current_table = {
                "pages": [page_number],
                "rows": []
            }
            tables.append(current_table)
        else:
            current_table["pages"].append(page_number)

        logger.debug(f"Contenuto: {len(chars)} caratteri, {len(rects)} celle, {len(images)} immagini.")

        previous_rect = None
        current_row = []

        previous_table_bottom = 0
        for rect in rects:
            if previous_rect is not None and rect["top"] != previous_rect["top"]:
                logger.debug("Trovato nuova riga")
                current_table["rows"].append({
                    "page": page_number,
                    "contents": current_row
                })
                current_row = []
            cell_chars = extract_cell_data(chars, rect)
            cell_text = join_chars(cell_chars)
            if "Numero domanda" == cell_text:
                logger.debug("Trovato inizio nuova tabella quiz")
                caption_text = ""
                if previous_table_bottom is not None:
                    caption_chars = extract_cell_data(chars, {
                        "x0": 0,
                        "x1": 600,
                        "top": previous_table_bottom,
                        "bottom": rect["top"]
                    })
                    caption_text = join_chars(caption_chars)
                    if bottom_text:
                        caption_text = f"{bottom_text} {caption_text}"
                        bottom_text = None
                if not caption_text:
                    logger.warning(f"Caption non recuperato: pag: {page_number}")
                else:
                    logger.info(f"Nuova tabella {caption_text}")
                current_table = {
                    "caption": caption_text,
                    "pages": [page_number],
                    "rows": []
                }
                tables.append(current_table)
            cell_image = extract_cell_data(images, rect)
            current_row.append({
                "text": cell_text,
                "images": [img["hash"] for img in cell_image]
            })
            previous_rect = rect
            previous_table_bottom = rect["bottom"]
        if current_row:
            current_table["rows"].append({
                    "page": page_number,
                    "contents": current_row
                })
        bottom_chars = extract_cell_data(chars, {
                        "x0": 0,
                        "x1": 600,
                        "top": previous_table_bottom,
                        "bottom": previous_table_bottom + 300
                    })
        if bottom_chars:
            bottom_text = join_chars(bottom_chars)
            logger.debug(f"Trovato testo a fine pagine fuori da tabelle {bottom_text}")
    logger.info(f"Estrazione tabelle completato con successo. Trovate {len(tables)} tabelle")
    return tables

        