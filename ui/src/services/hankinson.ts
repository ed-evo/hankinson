import { useLocalStorage } from "@vueuse/core";
import { ofetch } from "ofetch";
import { type Ref } from "vue";

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

export type User = string

export const USER_REF: Ref<User | null> = useLocalStorage('hankinson.user-email', null)

const baseURL = "/api/v1/quiz"

function default_request() {
    let headers: HeadersInit = {}
    if (USER_REF.value) {
        headers['X-Authenticated-User'] = USER_REF.value
    }
    return {
        baseURL,
        headers
    }
}

export async function login(): Promise<User> {
    try {
        const user = await ofetch<User>("/me", default_request())
        USER_REF.value = user
        return user
    } catch (err) {
        console.error("Quiz login error", err)
        throw err
    }
}

export async function getCapitoli(): Promise<Capitolo[]> {
    return ofetch<Capitolo[]>("/capitoli", default_request())
    
}

export async function getDomandeByCapitolo(capitoloId: number): Promise<Domanda[]> {
    const capitolo = await ofetch<Capitolo>(`/capitoli/${capitoloId}`, default_request())
    if (!capitolo?.edges?.domande) {
        throw new Error(`Domande non trovate per capitolo ${capitoloId}`)
    }
    return capitolo.edges.domande
}