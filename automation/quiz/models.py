from enum import StrEnum
from typing import Optional
from pydantic import BaseModel, TypeAdapter

class RispostaCorretta(StrEnum):
    VERO = "VERO"
    FALSO = "FALSO"

class Metadata(BaseModel):
    id_blocco: int
    argomenti: list[str]

class Domanda(Metadata):
    numero: int
    testo: str
    risposta_corretta: RispostaCorretta
    immagine: Optional[str] = None  # Giving optional fields a default of None
    pagina_quiz: int

DomandeList = TypeAdapter(list[Domanda])