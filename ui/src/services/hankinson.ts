import CapitoliSelect from "@/components/CapitoliSelect.vue";
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

export interface Quesito {
    id: number,
    esameId: number,
    domandaId: number 
}

export interface AttivitaQuesito {
    tipo: string,
    inizio: Date,
    fine: Date
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

function build_request_options(overrides: any = {}) {
    let headers: HeadersInit = {}
    if (USER_REF.value) {
        headers['X-Authenticated-User'] = USER_REF.value
    }
    return {
        baseURL,
        headers,
        ...overrides
    }
}

export async function login(): Promise<User> {
    try {
        const user = await ofetch<User>("/me", {
            ...build_request_options(),
            baseURL: '/api/v1'
        })
        USER_REF.value = user
        return user
    } catch (err) {
        console.error("Quiz login error", err)
        throw err
    }
}

export async function getCapitoli(): Promise<Capitolo[]> {
    return ofetch<Capitolo[]>("/capitoli", build_request_options())
    
}

export async function getDomandeByCapitolo(capitoloId: number): Promise<Domanda[]> {
    const capitolo = await ofetch<Capitolo>(`/capitoli/${capitoloId}`, build_request_options())
    if (!capitolo?.edges?.domande) {
        throw new Error(`Domande non trovate per capitolo ${capitoloId}`)
    }
    return capitolo.edges.domande
}

export async function getDomandaById(domandaId: number): Promise<Domanda> {
    return ofetch<Domanda>(`/domande/${domandaId}`, build_request_options())
}

export async function nextQuesitoAperto(capitoliIds: number[]): Promise<Quesito> {
    const quesito = await ofetch<Quesito>("/esami/aperto/next", build_request_options({
        method: 'POST',
        body: {
            capitoli: capitoliIds
        }
    }))
    return quesito
}

export async function notifyQuesityAttivita(quesitoId: number, attivita: AttivitaQuesito): Promise<void> {
    await ofetch(`/esami/quesiti/${quesitoId}/attivita`, build_request_options({
        method: 'PUT',
        body: attivita
    }))
}