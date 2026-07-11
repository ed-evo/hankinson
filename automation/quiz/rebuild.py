import json
from operator import itemgetter
from pathlib import Path

def extract_horizontal_lines(page, top, bottom):
    horizontal_lines = []
    for rect in page.rects:
        if rect["top"] >= top and rect["bottom"] <= bottom:
            horizontal_lines.append(rect)
    return horizontal_lines

def extract_vertical_lines(page, left, right):
    vertical_lines = []
    for rect in page.rects:
        if rect["x0"] >= left and rect["x1"] <= right:
            vertical_lines.append(rect)
    return vertical_lines

def extract_cell_data(page_data, cell_rect):
    return [
        item
            for item in page_data
            if cell_rect["x0"] <= item["x0"] <= cell_rect["x1"]
                and cell_rect["top"] <= item["top"] <= cell_rect["bottom"]
    ]

def extract_table_rect(rects):
    if not rects:
        return None

    min_x0 = min(rect["x0"] for rect in rects)
    max_x1 = max(rect["x1"] for rect in rects)
    min_top = min(rect["top"] for rect in rects)
    max_bottom = max(rect["bottom"] for rect in rects)

    return {
        "x0": min_x0,
        "x1": max_x1,
        "top": min_top,
        "bottom": max_bottom
    }

def join_chars(chars):
    if not chars:
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

    return " ".join(words)

def extract_quizes(pages):
    first_page = pages[0]
    header_rect = {
        "x0": 0,
        "x1": 600,
        "top": 0,
        "bottom": 100
    }
    first_page["rects"].insert(0, header_rect)

    by_top_x0 = itemgetter("top", "x0")

    tables = []
    current_table = None
    previous_table_bottom = None
    for page in pages:
        chars = sorted(page["chars"], key=by_top_x0)
        rects = sorted(page["rects"], key=by_top_x0)
        images = sorted(page["images"], key=by_top_x0)
        if current_table is None:
            current_table = {
                "pages": [page["pageNumber"]],
                "rows": []
            }
            tables.append(current_table)
        else:
            current_table["pages"].append(page["pageNumber"])

        print(f"📄 Page {page['pageNumber']} has {len(chars)} characters, {len(rects)} rectangles, and {len(images)} images.")

        previous_rect = None
        current_row = []
        for rect in rects:
            if previous_rect is not None and rect["top"] != previous_rect["top"]:
                current_table["rows"].append(current_row)
                current_row = []
            cell_chars = extract_cell_data(chars, rect)
            cell_text = join_chars(cell_chars)
            if "Numero domanda" == cell_text:
                caption_text = ""
                if previous_table_bottom is not None:
                    caption_chars = extract_cell_data(chars, {
                        "x0": 0,
                        "x1": 600,
                        "top": previous_table_bottom,
                        "bottom": rect["top"]
                    })
                    caption_text = join_chars(caption_chars)

                current_table = {
                    "caption": caption_text,
                    "pages": [page["pageNumber"]],
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
            current_table["rows"].append(current_row)
    return tables

if __name__ == "__main__":
    script_dir = Path(__file__).parent
    with open(script_dir / "pdf_geometry_dump.json", "r") as f:
        data = json.load(f)
    tables = extract_quizes(data)
    with open(script_dir / "quizes.json", "w") as f:
        json.dump(tables, f, indent=2)
        