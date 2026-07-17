
export interface Argomento {
    id: number;
    nome: string;
}
export interface Domanda {
    id: number;
    id_blocco: number;
    id_capitolo: number;
    immagine: string;
    is_true: boolean;
    pagina_quiz: number;
    testo: string;
    totale_domande: number;
    edges?: {
        argomenti?: Argomento[]
    }
}
export interface Capitolo {
    id: number;
    nome: string;
    max_numero_domanda: number;
    min_numero_domanda: number;
    totale_domande: number;
    edges?: {
        domande?: Domanda[]
    }
}