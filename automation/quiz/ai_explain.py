import io
import json
import os
from pathlib import Path
from pprint import pprint
from google import genai
from google.genai import types
from PIL import Image
from pydantic import BaseModel, Field, TypeAdapter
from typing import List
from dotenv import load_dotenv

from models import Domanda, DomandeList
load_dotenv()

class QuizExplanation(BaseModel):
    numero: int = Field(description="Il numero identificativo della domanda")
    spiegazione: str = Field(description="Spiegazione semplice e chiara del perché è VERO/FALSO")
    focus_linguistico: str = Field(description="Chiarimento sulle parole trappola (es. 'di norma')")
    regola_chiave: str = Field(description="Pillola mnemonica brevissima")

QuizExplanationList = TypeAdapter(list[QuizExplanation])

class BlockEnrichmentResponse(BaseModel):
    quiz_spiegati: List[QuizExplanation]

def generate_request(file_path: Path, assets_dir: Path) -> types.InlinedRequest:
    if not file_path.exists():
        raise FileNotFoundError(f"{file_path} Not found")

    domande: list[Domanda] = DomandeList.validate_json(file_path.read_bytes())

    image_path = None
    if domande[0].immagine:
        image_path = assets_dir / f"{domande[0].immagine}.png"


    system_instruction = (
        "Sei un insegnante di scuola guida. Ti verrà fornito un gruppo di domande legate a una stessa immagine o argomento.\n"
        "Per OGNI domanda presente nell'elenco, genera una spiegazione semplice (livello B1), evidenzia le parole trappola "
        "e fornisci una regola chiave mnemonica. Assicurati di includere il 'numero' identificativo per ciascuna domanda."
    )

    parts: list[types.Part] = []

    # Aggiungiamo l'immagine se disponibile
    if image_path and os.path.exists(image_path):
        parts.append(types.Part.from_bytes(
            data=open(image_path, 'rb').read(),
            mime_type='image/png'
        ))
    # Prepariamo il prompt con l'elenco di tutte le domande del blocco
    prompt_text = "Ecco le domande del blocco da analizzare:\n\n"
    for domanda in domande:
        prompt_text += f"- Numero: {domanda.numero}\n"
        prompt_text += f"  Testo: \"{domanda.testo}\"\n"
        prompt_text += f"  Risposta Corretta: {domanda.risposta_corretta}\n\n"

    parts.append(types.Part.from_text(text=prompt_text))

    content = types.UserContent(parts=parts)

    config = types.GenerateContentConfig(
        system_instruction=types.Content(parts=[types.Part.from_text(text=system_instruction)]),
        response_mime_type="application/json",
        response_schema=BlockEnrichmentResponse.model_json_schema(),
        temperature=0.2,
    )

    request = types.InlinedRequest(
        config=config,
        contents=[content]
    )

    return request

def read_batch_result_file(batch_job: types.BatchJob) -> list[QuizExplanation]:

    results: list[QuizExplanation] = []

    for response in batch_job.dest.inlined_responses:
        quiz_response = BlockEnrichmentResponse.model_validate_json(response.response.text)
        results.extend(quiz_response.quiz_spiegati)
    print(len(results))
    return results

# Esempio d'uso con il tuo JSON
if __name__ == "__main__":

    baseDir = Path(__file__).parent.resolve()
    buildDir = baseDir / 'build'
    assetsDir = buildDir / 'quiz_assets'
    batch_job_path = baseDir / 'batch.json'
    spiegazioni_out_path = buildDir / 'spiegazioni.json'

    spiegazioni = read_batch_result_file(
        types.BatchJob.model_validate_json(batch_job_path.read_text())
    )

    spiegazioni_out_path.write_bytes(QuizExplanationList.dump_json(spiegazioni, indent=2, exclude_none=True))


    # requests: list[types.InlinedRequest] = []
    # for file in buildDir.glob('quesito_*.json'):
    #     batch_request = generate_request(file, assetsDir)
    #     requests.append(batch_request)

    # client = genai.Client()


    # batch_job = client.batches.create(
    #     model='gemini-2.5-flash',
    #     src=requests,

    #     config={ 'display_name': 'quiz-job-full-by-blocco'}
    # )

    # print(f"Create batch job: {batch_job.name}")
    # print(batch_job.model_dump_json(exclude_none=True, indent=2))
    
    # print(Path(".").absolute())
    # with open("quiz/build/quesito_4001.json", "r", encoding="utf-8") as f:
    #     quiz_blocco = json.load(f)

    # # Percorso eventuale dell'immagine
    # img_id = quiz_blocco[0].get("immagine")
    # img_path = f"./quiz/build/quiz_assets/{img_id}.png"  # adatta l'estensione se .jpg

    # print(f"Inviando il blocco di {len(quiz_blocco)} domande in 1 sola richiesta...")
    # blocco_arricchito = enrich_block(quiz_blocco, image_path=img_path)

    # print(blocco_arricchito.model_dump_json(indent=2))

    # name = 'batches/93dr3rqyclc490bw0jtco01ipyhh6ngx83ax'
    # client = genai.Client()
    # batch_job = client.batches.get(name=name)

    # pprint(batch_job.dest.inlined_responses[0].response.text)
    # batch_job_path.open('w').write(batch_job.model_dump_json())

    # print("✅ Blocco completato e salvato!")
