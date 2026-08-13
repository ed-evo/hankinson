import type { Ref } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { ofetch } from 'ofetch'
import z from 'zod'

function booleanToRisposta (value?: boolean | null): RispostaEnum | null {
  if (value === undefined || value === null) {
    return null
  }
  return value ? RispostaEnum.VERO : RispostaEnum.FALSO
}

// ==========================================
// Enums
// ==========================================

export enum TipoEsameEnum {
  MINISTERIALE = 'ministeriale',
  PARZIALE = 'parziale',
  APERTO = 'aperto',
}

export enum TipoAttivitaEnum {
  SALTA = 'salta',
  RISPOSTA = 'risposta',
  PAUSA = 'pausa',
  PROSSIMO = 'prossimo',
}

export enum RispostaEnum {
  VERO = 'V',
  FALSO = 'F',
}

export enum TipoCorrezioneEnum {
  HUMAN = 'human',
  AI = 'ai',
}

export type User = string

export const USER_REF: Ref<User | null> = useLocalStorage(
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

export const CapitoloDtoSchema = z.object({
  id: z.number().int(),
  nome: z.string(),
  min_numero_domanda: z.number().int(),
  max_numero_domanda: z.number().int(),
  totale_domande: z.number().int(),
})

export const CapitoloSchema = CapitoloDtoSchema.transform(dto => ({
  id: dto.id,
  nome: dto.nome,
  minNumeroDomanda: dto.min_numero_domanda,
  maxNumeroDomanda: dto.max_numero_domanda,
  totaleDomande: dto.totale_domande,
}))

export type Capitolo = z.infer<typeof CapitoloSchema>

const CapitoliSchema = z.array(CapitoloSchema)

export async function getCapitoli (): Promise<Capitolo[]> {
  return hksApi<Capitolo[]>('/capitoli', {
    parseResponse: createParserFor(CapitoliSchema),
  })
}

export const CapitoloBasicStatsDtoSchema = z.object({
  id: z.number().int(),
  totale: z.number().int(),
  corrette: z.number().int(),
  sbagliate: z.number().int(),
  non_date: z.number().int(),
  durata_ms: z.number().int(),
})

export const CapitoloBasicStatsSchema = CapitoloBasicStatsDtoSchema.transform(
  dto => ({
    id: dto.id,
    totale: dto.totale,
    corrette: dto.corrette,
    sbagliate: dto.sbagliate,
    nonDate: dto.non_date,
    durataMs: dto.durata_ms,
  }),
)

export type CapitoloBasicStats = z.infer<typeof CapitoloBasicStatsSchema>

const CapitoliBasicStatsSchema = z.array(CapitoloBasicStatsSchema)
export async function getCapitoliStats (): Promise<CapitoloBasicStats[]> {
  return hksApi<CapitoloBasicStats[]>(
    '/capitoli/stats',
    {
      parseResponse: createParserFor(CapitoliBasicStatsSchema),
    },
  )
}

export const SpiegazioneDtoSchema = z.object({
  id: z.number().int().optional(),
  numero_domanda: z.number().int(),
  spiegazione: z.string(),
  focus_linguistico: z.string(),
  regola_chiave: z.string(),
})

export const SpiegazioneSchema = SpiegazioneDtoSchema.transform(dto => ({
  id: dto.id,
  numeroDomanda: dto.numero_domanda,
  spiegazione: dto.spiegazione,
  focusLinguistico: dto.focus_linguistico,
  regolaChiave: dto.regola_chiave,
}))

export type Spiegazione = z.infer<typeof SpiegazioneSchema>

export async function spiegaDomandaById (
  domandaId: number,
): Promise<Spiegazione> {
  return hksApi<Spiegazione>(
    `/domande/${domandaId}/spiegazione`,
    { method: 'POST', parseResponse: createParserFor(SpiegazioneSchema) },
  )
}

export const EsameDtoSchema = z.object({
  id: z.number().int(),
  tipo: z.enum(TipoEsameEnum),
  numero_quesiti: z.number().int(),
  max_errori: z.number().int(),
  minuti_disponibili: z.number().int(),
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
})

export const EsameSchema = EsameDtoSchema.transform(dto => {
  const erroriTotali = 0

  return {
    id: dto.id,
    tipo: dto.tipo,
    numeroQuesiti: dto.numero_quesiti,
    maxErrori: dto.max_errori,
    minutiDisponibili: dto.minuti_disponibili,
    createdAt: dto.created_at,
    updatedAt: dto.updated_at,
    erroriTotali,
    esitoSuperato: erroriTotali <= dto.max_errori,
  }
})

export type Esame = z.infer<typeof EsameSchema>

export async function getEsameById (esameId: number): Promise<Esame> {
  return await hksApi<Esame>(`/esami/${esameId}`, {
    parseResponse: createParserFor(EsameSchema),
  })
}

const EsameParzialeParamsInputSchema = z.object({
  capitoli: z.array(z.number().int()).min(1, 'Seleziona almeno un capitolo'),
  numeroQuesiti: z.number().int().positive(),
  maxErrori: z.number().int().nonnegative(),
  minutiDisponibili: z.number().int().positive(),
})

export type EsameParzialeParamsInput = z.infer<typeof EsameParzialeParamsInputSchema>

export const EsameParzialeParamsDtoSchema = EsameParzialeParamsInputSchema.transform(ui => ({
  capitoli: ui.capitoli,
  numero_quesiti: ui.numeroQuesiti,
  max_errori: ui.maxErrori,
  minuti_disponibili: ui.minutiDisponibili,
}))

export async function createEsameParziale (
  params: EsameParzialeParamsInput,
): Promise<Esame> {
  return hksApi<Esame>(
    '/esami/parziali',
    {
      method: 'PUT',
      body: EsameParzialeParamsDtoSchema.parse(params),
      parseResponse: createParserFor(EsameSchema),
    },
  )
}

export const CorrezioneDtoSchema = z.object({
  id: z.number().int(),
  esame_id: z.number().int(),
  type: z.enum(TipoCorrezioneEnum),
  esaminatore: z.string(),
  is_promosso: z.boolean(),
  testo: z.string(),
  meta: z.string(),
  created_at: z.coerce.date(),
})

export const CorrezioneSchema = CorrezioneDtoSchema.transform(dto => ({
  id: dto.id,
  esameId: dto.esame_id,
  type: dto.type,
  esaminatore: dto.esaminatore,
  isPromosso: dto.is_promosso,
  testo: dto.testo,
  meta: dto.meta,
  createdAt: dto.created_at,
}))

const CorrezioniSchema = z.array(CorrezioneSchema)

export type Correzione = z.infer<typeof CorrezioneSchema>

export async function getCorrezioniEsame (esameId: number): Promise<Correzione[]> {
  return hksApi<Correzione[]>(
    `/esami/${esameId}/correzioni`,
    {
      parseResponse: createParserFor(CorrezioniSchema),
    },
  )
}

export async function aiCorrege (esameId: number): Promise<Correzione[]> {
  return hksApi<Correzione[]>(
    `/esami/${esameId}/ai-corregge`,
    {
      method: 'POST',
      parseResponse: createParserFor(CorrezioniSchema),
    },
  )
}

export const QuesitoEsameDtoSchema = z.object({
  id: z.number().int(),
  risposta_finale: z.boolean().nullable().optional(),
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
})

const QuesitoEsameSchema = QuesitoEsameDtoSchema.transform(dto => {
  return {
    id: dto.id,
    rispostaFinale: booleanToRisposta(dto.risposta_finale),
    createdAt: dto.created_at,
    updatedAt: dto.updated_at,
  }
})

export type QuesitoEsame = z.infer<typeof QuesitoEsameSchema>

export async function nextQuesitoAperto (
  capitoliIds: number[],
): Promise<QuesitoEsame> {
  return await hksApi<QuesitoEsame>(
    '/esami/aperto/next',
    {
      method: 'POST',
      body: {
        capitoli: capitoliIds,
      },
      parseResponse: createParserFor(QuesitoEsameSchema),
    },
  )
}

const QuesitiSchema = z.array(QuesitoEsameSchema)
export async function getEsameQuesiti (
  esameId: number,
): Promise<QuesitoEsame[]> {
  return await hksApi<QuesitoEsame[]>(
    `/esami/${esameId}/quesiti`,
    {
      parseResponse: createParserFor(QuesitiSchema),
    },
  )
}

export const QuesitiBasicStatsDtoSchema = z.object({
  totale: z.number().int(),
  corrette: z.number().int(),
  sbagliate: z.number().int(),
  non_date: z.number().int(),
})

export const QuesitiBasicStatsSchema = QuesitiBasicStatsDtoSchema.transform(
  dto => ({
    totale: dto.totale,
    corrette: dto.corrette,
    sbagliate: dto.sbagliate,
    nonDate: dto.non_date,
  }),
)

export type QuesitiBasicStats = z.infer<typeof QuesitiBasicStatsSchema>

export async function getQuesitiStats (): Promise<QuesitiBasicStats> {
  return hksApi<QuesitiBasicStats>(
    '/esami/quesiti/stats',
    {
      parseResponse: createParserFor(QuesitiBasicStatsSchema),
    },
  )
}

export const DomandaDtoSchema = z.object({
  id: z.number().int(),
  testo: z.string(),
  is_true: z.boolean(),
  immagine: z.string().nullable().optional(),
  id_capitolo: z.number().int(),
  pagina_quiz: z.number().int(),
  id_blocco: z.number().int(),
})

export const DomandaSchema = DomandaDtoSchema.transform(dto => ({
  id: dto.id,
  testo: dto.testo,
  rispostaCorretta: dto.is_true ? RispostaEnum.VERO : RispostaEnum.FALSO,
  immagine: dto.immagine ?? null,
  immaginePath: dto.immagine ? `/quiz_assets/${dto.immagine}.png` : null,
  idCapitolo: dto.id_capitolo,
  paginaQuiz: dto.pagina_quiz,
  idBlocco: dto.id_blocco,
}))

export type Domanda = z.infer<typeof DomandaSchema>

export async function getQuestitoDomanda (quesitoId: number): Promise<Domanda> {
  return hksApi<Domanda>(
    `/esami/quesiti/${quesitoId}/domanda`,
    {
      parseResponse: createParserFor(DomandaSchema),
    },
  )
}

export const AttivitaQuesitoEsameDtoSchema = z.object({
  tipo: z.enum(TipoAttivitaEnum),
  risposta_data: z.boolean().nullable().optional(),
  inizio: z.coerce.date(),
  durata_ms: z.number().int(),
})

const AttivitaQuesitoEsameSchema = z.object({
  tipo: z.enum(TipoAttivitaEnum),
  rispostaData: z.enum(RispostaEnum).nullable().optional(),
  inizio: z.date(),
  durataMs: z.number().int(),
})

export type AttivitaQuesitoEsame = z.infer<typeof AttivitaQuesitoEsameSchema>

const AttivitaQuesitoCodec = z.codec(
  AttivitaQuesitoEsameDtoSchema,
  AttivitaQuesitoEsameSchema,
  {
    decode: dto => {
      return {
        tipo: dto.tipo,
        rispostaData: booleanToRisposta(dto.risposta_data),
        inizio: dto.inizio,
        durataMs: dto.durata_ms,
      }
    },
    encode: model => {
      return {
        tipo: model.tipo,
        risposta_data: model.rispostaData ? model.rispostaData === RispostaEnum.VERO : undefined,
        inizio: model.inizio,
        durata_ms: model.durataMs,
      }
    },
  },
)

export async function listQuesitoAttivita (quesitoId: number): Promise<AttivitaQuesitoEsame[]> {
  return await hksApi<AttivitaQuesitoEsame[]>(
    `/esami/quesiti/${quesitoId}/attivita`,
    {
      parseResponse (responseText: string) {
        return JSON.parse(responseText).map((a: any) => AttivitaQuesitoCodec.decode(a))
      },
    },
  )
}

export async function notifyQuesityAttivita (
  quesitoId: number,
  attivita: AttivitaQuesitoEsame,
): Promise<void> {
  await hksApi(
    `/esami/quesiti/${quesitoId}/attivita`,
    {
      method: 'PUT',
      body: AttivitaQuesitoCodec.encode(attivita),
    },
  )
}

export interface PausaEvent {
  inizio: Date
  fine: Date
}
export class AttivitaEmitter {
  constructor (private startedAt: Date = new Date()) {}

  reset (at: Date = new Date()) {
    this.startedAt = at
  }

  async fire (
    idQuesito: number,
    tipo: TipoAttivitaEnum,
    risposta: RispostaEnum | null = null,
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

  async firePausa (idQuesito: number, event: PausaEvent) {
    await notifyQuesityAttivita(idQuesito, {
      tipo: TipoAttivitaEnum.PAUSA,
      inizio: event.inizio,
      durataMs: event.fine.getTime() - event.inizio.getTime(),
    })
  }
}

// unused scheme

export const ArgomentoDtoSchema = z.object({
  id: z.number().int(),
  nome: z.string(),
})

export const ArgomentoSchema = ArgomentoDtoSchema.transform(dto => ({
  id: dto.id,
  nome: dto.nome,
}))

export type Argomento = z.infer<typeof ArgomentoSchema>
