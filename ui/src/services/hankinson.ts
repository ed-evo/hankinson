import type { Ref } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { ofetch } from 'ofetch'

export enum Choice {
  VERO = 'V',
  FALSO = 'F',
}
export interface Argomento {
  id: number
  nome: string
}
export interface Domanda {
  id: number
  id_blocco: number
  id_capitolo: number
  immagine: string
  is_true: boolean
  pagina_quiz: number
  testo: string
  totale_domande: number
  edges?: {
    argomenti?: Argomento[]
  }
}

export interface SpiegazioneDomanda {
  id: number
  spiegazione: string
  focus_linguistico: string
  regola_chiave: string
}

export interface QuesitoEsame {
  id: number
  risposta_finale?: boolean
  edges: {
    domanda_originale: Domanda
    attivita: AttivitaQuesito[]
  }
}

export interface Esame {
  id: number
  tipo: string
  max_errori: number
  numero_quesiti: number
  minuti_disponibili: number
  edges: {
    quesiti: QuesitoEsame[]
  }
}

export interface EsameParzialeParams {
  capitoli: number[]
  numero_quesiti: number
  max_errori: number
  minuti_disponibili: number
}

export interface Quesito {
  id: number
  esameId: number
  domandaId: number
}

export interface QuesitiBasicStats {
  totale: number
  corrette: number
  sbagliate: number
  non_date: number
}

export enum TipoAttivitaQuesito {
  salta = 'salta',
  risposta = 'risposta',
  pausa = 'pausa',
  prossimo = 'prossimo',
}

export interface AttivitaQuesito {
  tipo: TipoAttivitaQuesito
  risposta_data?: boolean
  inizio: Date
  durata_ms: number
}

export interface PausaEvent {
  inizio: Date
  fine: Date
}

export interface Capitolo {
  id: number
  nome: string
  max_numero_domanda: number
  min_numero_domanda: number
  totale_domande: number
  edges?: {
    domande?: Domanda[]
  }
}

export interface CapitoloBasicStats {
  id: number
  totale: number
  corrette: number
  sbagliate: number
  non_date: number
  durata_ms: number
}

export type User = string

export function getImmaginePath(domanda: Domanda): string | undefined {
  if (domanda?.immagine) {
    return `/quiz_assets/${domanda.immagine}.png`
  }
}

export const USER_REF: Ref<User | null> = useLocalStorage(
  'hankinson.user-email',
  null,
)

const baseURL = '/api/v1/quiz'

const hksApi = ofetch.create({
  baseURL,
  onRequest({ options }) {

    if (USER_REF.value) {
      options.headers.set('X-Authenticated-User', USER_REF.value);
    }
  }
})

export async function login (): Promise<User> {
  try {
    const user = await hksApi<User>('/me', {
      baseURL: '/api/v1',
    })
    USER_REF.value = user
    return user
  } catch (error) {
    console.error('Quiz login error', error)
    throw error
  }
}

export async function getCapitoli (): Promise<Capitolo[]> {
  return hksApi<Capitolo[]>('/capitoli')
}

export async function getCapitoliStats (): Promise<CapitoloBasicStats[]> {
  return hksApi<CapitoloBasicStats[]>(
    '/capitoli/stats',
  )
}

export async function getDomandeByCapitolo (
  capitoloId: number,
): Promise<Domanda[]> {
  const capitolo = await hksApi<Capitolo>(
    `/capitoli/${capitoloId}`,
  )
  if (!capitolo?.edges?.domande) {
    throw new Error(`Domande non trovate per capitolo ${capitoloId}`)
  }
  return capitolo.edges.domande
}

export async function getDomandaById (domandaId: number): Promise<Domanda> {
  return hksApi<Domanda>(`/domande/${domandaId}`)
}

export async function spiegaDomandaById (
  domandaId: number,
): Promise<SpiegazioneDomanda> {
  return hksApi<SpiegazioneDomanda>(
    `/domande/${domandaId}/spiegazione`,
   { method: 'POST' },
  )
}

export async function nextQuesitoAperto (
  capitoliIds: number[],
): Promise<Quesito> {
  const quesito = await hksApi<Quesito>(
    '/esami/aperto/next',
    {
      method: 'POST',
      body: {
        capitoli: capitoliIds,
      },
    },
  )
  return quesito
}

export async function getQuesitiStats (): Promise<QuesitiBasicStats> {
  return hksApi<QuesitiBasicStats>(
    '/esami/quesiti/stats',
  )
}

export async function notifyQuesityAttivita (
  quesitoId: number,
  attivita: AttivitaQuesito,
): Promise<void> {
  await hksApi(
    `/esami/quesiti/${quesitoId}/attivita`,
    {
      method: 'PUT',
      body: attivita,
    },
  )
}

export async function getEsameQuesiti (
  esameId: number,
): Promise<QuesitoEsame[]> {
  return await hksApi<QuesitoEsame[]>(
    `/esami/${esameId}/quesiti`,
  )
}

export async function getEsameById (esameId: number): Promise<Esame> {
  return await hksApi<Esame>(`/esami/${esameId}`)
}

export async function createEsameParziale (
  params: EsameParzialeParams,
): Promise<Esame> {
  return hksApi<Esame>(
    '/esami/parziali',
    {
      method: 'PUT',
      body: params,
    },
  )
}

export class AttivitaEmitter {
  constructor (private startedAt: Date = new Date()) {}

  reset (at: Date = new Date()) {
    this.startedAt = at
  }

  async fire (
    idQuesito: number,
    tipo: TipoAttivitaQuesito,
    risposta: Choice | null = null,
    finishedAt: Date = new Date(),
  ) {
    let risposta_data = undefined
    if (risposta) {
      risposta_data = Choice.VERO === risposta
    }
    await notifyQuesityAttivita(idQuesito, {
      tipo,
      inizio: this.startedAt,
      durata_ms: finishedAt.getTime() - this.startedAt.getTime(),
      risposta_data,
    })
    this.startedAt = finishedAt
  }

  async firePausa (idQuesito: number, event: PausaEvent) {
    await notifyQuesityAttivita(idQuesito, {
      tipo: TipoAttivitaQuesito.pausa,
      inizio: event.inizio,
      durata_ms: event.fine.getTime() - event.inizio.getTime(),
    })
  }
}
