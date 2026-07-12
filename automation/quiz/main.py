import json
import argparse
from pathlib import Path
import shutil
from collections import defaultdict

import context
import extract
import rebuild
import normalize

pdf_name = "DomandeB.pdf"
build_dirname = "build"
assets_dirname = "quiz_assets"

if __name__ == "__main__":
    # Set up the argument parser
    parser = argparse.ArgumentParser(
        description="Estrazione Quiz per patente B",
        allow_abbrev=False
        )

    parser.add_argument(
        '--log', 
        default='INFO', 
        choices=['DEBUG', 'INFO', 'WARNING', 'ERROR'],
        help="Set log level"
    )
    
    # Add the intermediate flag (True if present, False if absent)
    parser.add_argument(
        '--intermediate', 
        action='store_true', 
        help="Salva risultati intermedi"
    )

    parser.add_argument(
        '--dry',
        action='store_true',
        help="Dry run"
    )

    parser.add_argument(
        '--offset',
        type=int,
        default=0,
        help="Numero di elementi da saltare all'inizio (default: 0)"
    )

    parser.add_argument(
        '--take',
        type=int,
        default=None,  # None means "take all remaining items" unless specified
        help="Numero massimo di elementi da estrarre (default: tutti)"
    )
    # Parse the arguments
    args = parser.parse_args()

    dry_run = args.dry

    context.DRY_RUN = dry_run

    context.LOG_LEVEL = args.log

    logger = context.getLogger("quiz-patente-b")

    save_intermediate = args.intermediate

    cwd = Path(__file__).parent
    pdf_path = cwd / pdf_name

    output_dir = cwd / build_dirname

    assets_dir = output_dir / assets_dirname
        
    logger.info("===================================")
    logger.info("Estrazione Quiz per patente B")
    logger.info(f"  Sorgente: {pdf_name}")
    logger.info(f"  Cartella di lavoro: {cwd}")
    logger.info(f"  Cartella risultati: {output_dir}")
    logger.info(f"  Cartella immagini: {assets_dir}")
    logger.info(f"  Salvataggio risultati intermedi: {'SI' if save_intermediate else 'NO'}")
    logger.info("===================================")

    if not dry_run:
        shutil.rmtree(output_dir)
        assets_dir.mkdir(exist_ok=True, parents=True)
    
    logger.info("Preparato cartella di output.")

    raw_data = extract.extract_pdf_with_clean_assets(
        pdf_path, assets_dir,
        offset=args.offset, take=args.take
    )

    if save_intermediate:
        if context.DRY_RUN:
            logger.info(json.dumps(raw_data, indent=2))
        else:
            intermedate_path = output_dir / "quiz_raw_data.json"
            with intermedate_path.open("w") as intermediate_file:
                json.dump(raw_data, intermediate_file, indent=2, ensure_ascii=False)
    
    # first table is header,
    # discarded to prevent further handling
    raw_tables = rebuild.extract_quizes(raw_data)[1:]
    if save_intermediate:
        if context.DRY_RUN:
            logger.info(json.dumps(raw_tables, indent=2))
        else:
            intermedate_path = output_dir / "quiz_raw_table.json"
            with intermedate_path.open("w") as intermediate_file:
                json.dump(raw_tables, intermediate_file, indent=2, ensure_ascii=False)
    
    domande = normalize.extract_quesiti(raw_tables)
    groupped = defaultdict(list)

    for domanda in domande:
        groupped[domanda["id_blocco"]].append(domanda)

    if context.DRY_RUN:
        logger.info(json.dumps(groupped, indent=2))
    else:
        for blocco, domande in groupped.items():
            filename = output_dir / f"quesito_{blocco}.json"
            with filename.open("w") as jsonfile:
                json.dump(domande, jsonfile, indent=2, ensure_ascii=False)
            logger.info(f"Quesito {blocco} con {len(domande)} domande salvato su {filename}")

    logger.info("Estrazione Quiz Patente B terminato correttamente.")