# /// script
# dependencies = [
#   "pdfplumber",
# ]
# ///

import hashlib
import io
import json
from pathlib import Path
import pdfplumber
from PIL import Image

def save_clean_image(page_obj, img_dict, output_dir: Path) -> str:
    """
    Uses Pillow to parse the raw byte stream directly out of the PDF object,
    avoiding heavy canvas rendering completely.
    """
    try:
        # 1. Grab the raw compressed stream data dictionary from pdfminer
        stream_data = img_dict.get("stream")
        if not stream_data:
            return "unknown_hash"
            
        raw_bytes = stream_data.get_data()
        if not raw_bytes:
            return "unknown_hash"

        # 2. Compute the MD5 hash directly from the raw bytes to deduplicate
        img_hash = hashlib.md5(raw_bytes).hexdigest()
        
        # Check if we already handled this specific traffic sign asset
        target_path = output_dir / f"{img_hash}.png"
        if target_path.exists():
            return img_hash

        # 3. Read metadata required by Pillow to form the image grid matrix
        width = img_dict.get("width")
        height = img_dict.get("height")
        colorspace = img_dict.get("colorspace", ["DeviceRGB"])
        
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
            img = Image.frombytes(mode, (int(width), int(height)), raw_bytes)
            img.save(target_path, format="PNG")
        return img_hash
        
    except Exception as e:
        print(f"⚠️ Pillow failed to reconstruct stream data: {e}")
        return "unknown_hash"

def extract_pdf_with_clean_assets(pdf_name="DomandeB.pdf", json_name="pdf_geometry_dump.json", assets_dirname="quiz_assets"):
    script_dir = Path(__file__).parent
    pdf_path = script_dir / pdf_name
    json_path = script_dir / json_name
    assets_dir = script_dir / assets_dirname

    assets_dir.mkdir(exist_ok=True)

    if not pdf_path.exists():
        print(f"❌ Error: {pdf_name} not found.")
        return

    extracted_document = []

    print(f"🚀 Capturing vector geometry layout and exporting assets to ./{assets_dirname}...")
    with pdfplumber.open(pdf_path) as pdf:
        for page in pdf.pages:
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

            # Map row grids and boxes
            for rect in page.rects:
                page_data["rects"].append({
                    "x0": round(rect["x0"], 2),
                    "x1": round(rect["x1"], 2),
                    "top": round(rect["top"], 2),
                    "bottom": round(rect["bottom"], 2)
                })

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

            extracted_document.append(page_data)
            print(f"  Page {page.page_number} geometry mapped completely.")

    # Write out structural data
    with open(json_path, "w", encoding="utf-8") as f:
        json.dump(extracted_document, f, indent=2, ensure_ascii=False)

    print(f"\n✅ Finished! All images are now standard, readable PNGs.")
    print(f"📂 Open look inside your asset folder: {assets_dir}")

if __name__ == "__main__":
    extract_pdf_with_clean_assets()