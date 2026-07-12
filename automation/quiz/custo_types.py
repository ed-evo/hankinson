from typing import Optional, TypedDict
from enum import StrEnum

class Metadata(TypedDict):
    id_blocco: int
    argomenti: list[str]

class RispostaCorretta(StrEnum):
    VERO = "VERO"
    FALSO = "FALSO"

class Domanda(Metadata):
    numero: int
    testo: str
    risposta_corretta: RispostaCorretta
    immagine: Optional[str]
    pagina_quiz: int