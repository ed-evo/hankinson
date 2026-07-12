import hashlib
import io
from pathlib import Path
import pdfplumber
from PIL import Image

import context

def save_clean_image(page_obj, img_dict, output_dir: Path) -> str:
    """
    Uses Pillow to parse the raw byte stream directly out of the PDF object,
    avoiding heavy canvas rendering completely.
    """
    logger = context.getLogger()
    try:
        # 1. Grab the raw compressed stream data dictionary from pdfminer
        stream_data = img_dict.get("stream")
        if not stream_data:
            logger.error("Trovato elemento immagine privo di stream")
            raise AssertionError("Immagine senza dati")
            
        raw_bytes = stream_data.get_data()
        if not raw_bytes:
            logger.error("Stream dell'immagine privo di dati")
            raise AssertionError("Immagine senza dati")

        # 2. Compute the MD5 hash directly from the raw bytes to deduplicate
        img_hash = hashlib.md5(raw_bytes).hexdigest()
        logger.debug(f"Trovato immagine: {img_dict.get('name')} - {img_hash}")

        if context.DRY_RUN:
            return img_hash
        
        # Check if we already handled this specific traffic sign asset
        target_path = output_dir / f"{img_hash}.png"
        if target_path.exists():
            logger.debug("Immagine esiste")
            return img_hash

        # 3. Read metadata required by Pillow to form the image grid matrix
        width = img_dict.get("width")
        height = img_dict.get("height")
        colorspace = img_dict.get("colorspace", ["DeviceRGB"])

        logger.debug(f"Immagine {width}x{height} {colorspace}")
        
        # Determine the color mode string for Pillow (e.g., RGB, L for Grayscale, CMYK)
        mode = "RGB"
        if "DeviceGray" in str(colorspace):
            mode = "L"
        elif "DeviceCMYK" in str(colorspace):
            mode = "CMYK"

        # 4. Attempt to pass the raw stream directly to Pillow
        # If the PDF stream is a native JPEG (DCTDecode), Pillow's Image.open can read the stream bytes directly
        try:
            img = Image.open(io.BytesIO(raw_bytes))
            img.save(target_path, format="PNG")
        except Exception:
            # If Image.open fails, it means the stream is un-headered raw pixel data (e.g., FlateDecode)
            # We reconstruct it manually using Pillow's frombytes factory mapping the metadata
            logger.info("Immagine 'un-headered raw pixel data")
            img = Image.frombytes(mode, (int(width), int(height)), raw_bytes)
            img.save(target_path, format="PNG")
        logger.debug(f"Salveto: {target_path}")
        return img_hash
        
    except Exception as e:
        logger.error(f"⚠️ Pillow failed to reconstruct stream data: {e}")
        raise e

def get_page_safe(items, offset: int = 0, take: int = None):
    # Prevent negative values crashing the logic
    offset = max(0, offset)
    
    if take is None:
        return items[offset:] # Take everything left if no 'take' is given
        
    take = max(0, take)
    return items[offset : offset + take]

def extract_pdf_with_clean_assets(pdf_path: Path, assets_dir: Path, offset: int = 0, take: int = None):
    logger = context.getLogger()

    if not pdf_path.exists():
        logger.error(f"❌ Errore: {pdf_path} non trovato.")
        raise FileNotFoundError(pdf_path)

    extracted_document = []

    with pdfplumber.open(pdf_path) as pdf:
        logger.info("Inizio estrazione quiz")
        logger.debug("Pdf aperto correttamente")
        logger.debug(f"Documento con {len(pdf.pages)}")
        # for page in pdf.pages:
        for page in get_page_safe(pdf.pages, offset=offset, take=take):
            logger.info(f"Pagina: {page.page_number}")
            page_data = {
                "pageNumber": page.page_number,
                "chars": [],
                "rects": [],
                "images": []
            }

            # Map text characters using standard Latin positions
            for char in page.chars:
                page_data["chars"].append({
                    "x0": round(char["x0"], 2),
                    "x1": round(char["x1"], 2),
                    "top": round(char["top"], 2),
                    "bottom": round(char["bottom"], 2),
                    "text": char["text"]
                })
            logger.info(f"Trovati {len(page_data['chars'])} caratteri.")

            # Map row grids and boxes
            for rect in page.rects:
                page_data["rects"].append({
                    "x0": round(rect["x0"], 2),
                    "x1": round(rect["x1"], 2),
                    "top": round(rect["top"], 2),
                    "bottom": round(rect["bottom"], 2)
                })

            logger.info(f"Trovate {len(page_data['rects'])} celle")

            # Target image dictionaries and pass the containing page object along
            for img in page.images:
                img_hash = save_clean_image(page, img, assets_dir)
                page_data["images"].append({
                    "x0": round(img["x0"], 2),
                    "x1": round(img["x1"], 2),
                    "top": round(img["top"], 2),
                    "bottom": round(img["bottom"], 2),
                    "name": img.get("name", "unnamed"),
                    "hash": img_hash
                })

            logger.info(f"Trovate {len(page_data['images'])} immagini")

            extracted_document.append(page_data)
            logger.info(f"Dati estratti correttamente dalla pagina {page.page_number}")


    logger.info(f"✅ Dati estratti correttamente")
    logger.info(f"📂 Immagini salvati nella cartella: {assets_dir}")
    return extracted_document
