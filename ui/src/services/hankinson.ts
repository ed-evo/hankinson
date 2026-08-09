import type { Ref } from 'vue'
import z from 'zod'
import { useLocalStorage } from '@vueuse/core'
import { ofetch } from 'ofetch'
import * as schema from '@/types/hankinson'

export const USER_REF: Ref<schema.User | null> = useLocalStorage(
  'hankinson.user-email',
  null,
)

const baseURL = '/api/v1/quiz'

const hksApi = ofetch.create({
  baseURL,
  onRequest ({ options }) {
    if (USER_REF.value) {
      options.headers.set('X-Authenticated-User', USER_REF.value)
    }
  },
})

function createParserFor<TOutput> (zSchema: z.ZodType<TOutput>): (responseText: string) => TOutput {
  return (responseText: string) => zSchema.parse(JSON.parse(responseText))
}

export async function login (): Promise<schema.User> {
  try {
    const user = await hksApi<schema.User>('/me', {
      baseURL: '/api/v1',
    })
    USER_REF.value = user
    return user
  } catch (error) {
    console.error('Quiz login error', error)
    throw error
  }
}

const CapitoliSchema = z.array(schema.CapitoloSchema)

export async function getCapitoli (): Promise<schema.Capitolo[]> {
  return hksApi<schema.Capitolo[]>('/capitoli', {
    parseResponse: createParserFor(CapitoliSchema),
  })
}

const CapitoliBasicStatsSchema = z.array(schema.CapitoloBasicStatsSchema)
export async function getCapitoliStats (): Promise<schema.CapitoloBasicStats[]> {
  return hksApi<schema.CapitoloBasicStats[]>(
    '/capitoli/stats',
    {
      parseResponse: createParserFor(CapitoliBasicStatsSchema),
    },
  )
}

const DomandeSchema = z.array(schema.CapitoloSchema)
export async function getDomandeByCapitolo (
  capitoloId: number,
): Promise<schema.Domanda[]> {
  const capitolo = await hksApi<schema.Capitolo>(
    `/capitoli/${capitoloId}`,
    {
      parseResponse: createParserFor(DomandeSchema),
    },
  )
  if (!capitolo?.domande) {
    throw new Error(`Domande non trovate per capitolo ${capitoloId}`)
  }
  return capitolo.domande
}

export async function getDomandaById (domandaId: number): Promise<schema.Domanda> {
  return hksApi<schema.Domanda>(`/domande/${domandaId}`, {
    parseResponse: createParserFor(schema.DomandaSchema),
  })
}

export async function spiegaDomandaById (
  domandaId: number,
): Promise<schema.Spiegazione> {
  return hksApi<schema.Spiegazione>(
    `/domande/${domandaId}/spiegazione`,
    { method: 'POST', parseResponse: createParserFor(schema.SpiegazioneSchema) },
  )
}

export async function nextQuesitoAperto (
  capitoliIds: number[],
): Promise<schema.QuesitoEsame> {
  return await hksApi<schema.QuesitoEsame>(
    '/esami/aperto/next',
    {
      method: 'POST',
      body: {
        capitoli: capitoliIds,
      },
      parseResponse: createParserFor(schema.QuesitoEsameSchema),
    },
  )
}

export async function getQuesitiStats (): Promise<schema.QuesitiBasicStats> {
  return hksApi<schema.QuesitiBasicStats>(
    '/esami/quesiti/stats',
    {
      parseResponse: createParserFor(schema.QuesitiBasicStatsSchema),
    },
  )
}

export async function notifyQuesityAttivita (
  quesitoId: number,
  attivita: schema.AttivitaQuesitoEsame,
): Promise<void> {
  let risposta = undefined
  if (attivita.rispostaData) {
    risposta = attivita.rispostaData == schema.RispostaEnum.VERO
  }
  await hksApi(
    `/esami/quesiti/${quesitoId}/attivita`,
    {
      method: 'PUT',
      body: schema.AttivitaQuesitoEsameDtoSchema.parse({
        tipo: attivita.tipo,
        risposta_data: risposta,
        inizio: attivita.inizio.toISOString(),
        durata_ms: attivita.durataMs,
      } as schema.AttivitaQuesitoEsameDto),
    },
  )
}

const QuesitiSchema = z.array(schema.QuesitoEsameSchema)
export async function getEsameQuesiti (
  esameId: number,
): Promise<schema.QuesitoEsame[]> {
  return await hksApi<schema.QuesitoEsame[]>(
    `/esami/${esameId}/quesiti`,
    {
      parseResponse: createParserFor(QuesitiSchema),
    },
  )
}

export async function getEsameById (esameId: number): Promise<schema.Esame> {
  return await hksApi<schema.Esame>(`/esami/${esameId}`, {
    parseResponse: createParserFor(schema.EsameSchema),
  })
}

export async function createEsameParziale (
  params: schema.EsameParzialeParamsInput,
): Promise<schema.Esame> {
  return hksApi<schema.Esame>(
    '/esami/parziali',
    {
      method: 'PUT',
      body: schema.EsameParzialeParamsDtoSchema.parse(params),
      parseResponse: createParserFor(schema.EsameSchema),
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
    tipo: schema.TipoAttivitaEnum,
    risposta: schema.RispostaEnum | null = null,
    finishedAt: Date = new Date(),
  ) {
    await notifyQuesityAttivita(idQuesito, {
      tipo,
      inizio: this.startedAt,
      durataMs: finishedAt.getTime() - this.startedAt.getTime(),
      rispostaData: risposta,
    })
    this.startedAt = finishedAt
  }

  async firePausa (idQuesito: number, event: schema.PausaEvent) {
    await notifyQuesityAttivita(idQuesito, {
      tipo: schema.TipoAttivitaEnum.PAUSA,
      inizio: event.inizio,
      durataMs: event.fine.getTime() - event.inizio.getTime(),
    })
  }
}
