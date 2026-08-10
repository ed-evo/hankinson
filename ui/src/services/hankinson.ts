import type { Ref } from 'vue'
import z from 'zod'
import { useLocalStorage } from '@vueuse/core'
import { ofetch } from 'ofetch'

function booleanToRisposta (value?: boolean | null): RispostaEnum | null {
  if (value === undefined || value === null) {
    return null
  }
  return value ? RispostaEnum.VERO : RispostaEnum.FALSO
}

// ==========================================
// 1. Dto
// ==========================================

// -- Enums

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

// --- DTO Types (API Backend Go) ---

type SpiegazioneDto = {
  id?: number
  numero_domanda: number
  spiegazione: string
  focus_linguistico: string
  regola_chiave: string
  edges?: {
    domanda?: DomandaDto
  }
}

type ArgomentoDto = {
  id: number
  nome: string
  edges?: {
    domande?: DomandaDto[]
  }
}

type CapitoloDto = {
  id: number
  nome: string
  min_numero_domanda: number
  max_numero_domanda: number
  totale_domande: number
  edges?: {
    domande?: DomandaDto[]
  }
}

type CapitoloBasicStatsDto = {
  id: number
  totale: number
  corrette: number
  sbagliate: number
  non_date: number
  durata_ms: number
}

type DomandaDto = {
  id: number
  testo: string
  is_true: boolean
  immagine?: string | null
  id_capitolo: number
  pagina_quiz: number
  id_blocco: number
  edges?: {
    argomenti?: ArgomentoDto[]
    capitolo?: CapitoloDto
    spiegazione?: SpiegazioneDto[]
  }
}

export type AttivitaQuesitoEsameDto = {
  id?: number
  tipo: TipoAttivitaEnum
  risposta_data?: boolean | null
  inizio: Date
  durata_ms: number
  timestamp?: Date
  edges?: {
    quesito_esame?: QuesitoEsameDto
  }
}

type QuesitoEsameDto = {
  id?: number
  risposta_finale?: boolean | null
  created_at: Date
  updated_at: Date
  edges?: {
    esame?: EsameDto
    domanda_originale?: DomandaDto
    attivita?: AttivitaQuesitoEsameDto[]
  }
}

type QuesitiBasicStatsDto = {
  totale: number
  corrette: number
  sbagliate: number
  non_date: number
}

type EsameDto = {
  id?: number
  tipo: TipoEsameEnum
  numero_quesiti: number
  max_errori: number
  minuti_disponibili: number
  created_at: Date
  updated_at: Date
  edges?: {
    utente?: any
    quesiti?: QuesitoEsameDto[]
  }
}

export type EsameParzialeParamsDto = {
  capitoli: number[]
  numero_quesiti: number
  max_errori: number
  minuti_disponibili: number
}

// ==========================================
// 2. Schemi DTO (con Getters per i campi ciclici)
// ==========================================

export const SpiegazioneDtoSchema: z.ZodType<SpiegazioneDto> = z.object({
  id: z.number().int().optional(),
  numero_domanda: z.number().int(),
  spiegazione: z.string(),
  focus_linguistico: z.string(),
  regola_chiave: z.string(),
  edges: z.object({
    get domanda (): z.ZodOptional<z.ZodType<DomandaDto>> {
      return DomandaDtoSchema.optional()
    },
  }).optional(),
})

export const ArgomentoDtoSchema: z.ZodType<ArgomentoDto> = z.object({
  id: z.number().int(),
  nome: z.string(),
  edges: z.object({
    get domande (): z.ZodOptional<z.ZodArray<z.ZodType<DomandaDto>>> {
      return z.array(DomandaDtoSchema).optional()
    },
  }).optional(),
})

export const CapitoloDtoSchema: z.ZodType<CapitoloDto> = z.object({
  id: z.number().int(),
  nome: z.string(),
  min_numero_domanda: z.number().int(),
  max_numero_domanda: z.number().int(),
  totale_domande: z.number().int(),
  edges: z.object({
    get domande (): z.ZodOptional<z.ZodArray<z.ZodType<DomandaDto>>> {
      return z.array(DomandaDtoSchema).optional()
    },
  }).optional(),
})

export const CapitoloBasicStatsDtoSchema: z.ZodType<CapitoloBasicStatsDto> = z.object({
  id: z.number().int(),
  totale: z.number().int(),
  corrette: z.number().int(),
  sbagliate: z.number().int(),
  non_date: z.number().int(),
  durata_ms: z.number().int(),
})

export const DomandaDtoSchema: z.ZodType<DomandaDto> = z.object({
  id: z.number().int(),
  testo: z.string(),
  is_true: z.boolean(),
  immagine: z.string().nullable().optional(),
  id_capitolo: z.number().int(),
  pagina_quiz: z.number().int(),
  id_blocco: z.number().int(),
  edges: z.object({
    get argomenti (): z.ZodOptional<z.ZodArray<z.ZodType<ArgomentoDto>>> {
      return z.array(ArgomentoDtoSchema).optional()
    },
    get capitolo (): z.ZodOptional<z.ZodType<CapitoloDto>> {
      return CapitoloDtoSchema.optional()
    },
    get spiegazione (): z.ZodOptional<z.ZodArray<z.ZodType<SpiegazioneDto>>> {
      return z.array(SpiegazioneDtoSchema).optional()
    },
  }).optional(),
})

export const AttivitaQuesitoEsameDtoSchema: z.ZodType<AttivitaQuesitoEsameDto> = z.object({
  id: z.number().int().optional(),
  tipo: z.enum(TipoAttivitaEnum),
  risposta_data: z.boolean().nullable().optional(),
  inizio: z.coerce.date(),
  durata_ms: z.number().int(),
  timestamp: z.coerce.date().optional(),
  edges: z.object({
    get quesito_esame (): z.ZodOptional<z.ZodType<QuesitoEsameDto>> {
      return QuesitoEsameDtoSchema.optional()
    },
  }).optional(),
})

export const QuesitoEsameDtoSchema: z.ZodType<QuesitoEsameDto> = z.object({
  id: z.number().int().optional(),
  risposta_finale: z.boolean().nullable().optional(),
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
  edges: z.object({
    get esame (): z.ZodOptional<z.ZodType<EsameDto>> {
      return EsameDtoSchema.optional()
    },
    get domanda_originale (): z.ZodOptional<z.ZodType<DomandaDto>> {
      return DomandaDtoSchema.optional()
    },
    get attivita (): z.ZodOptional<z.ZodArray<z.ZodType<AttivitaQuesitoEsameDto>>> {
      return z.array(AttivitaQuesitoEsameDtoSchema).optional()
    },
  }).optional(),
})

export const QuesitiBasicStatsDtoSchema: z.ZodType<QuesitiBasicStatsDto> = z.object({
  totale: z.number().int(),
  corrette: z.number().int(),
  sbagliate: z.number().int(),
  non_date: z.number().int(),
})

export const EsameDtoSchema: z.ZodType<EsameDto> = z.object({
  id: z.number().int().optional(),
  tipo: z.enum(TipoEsameEnum),
  numero_quesiti: z.number().int(),
  max_errori: z.number().int(),
  minuti_disponibili: z.number().int(),
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
  edges: z.object({
    utente: z.any().optional(),
    get quesiti (): z.ZodOptional<z.ZodArray<z.ZodType<QuesitoEsameDto>>> {
      return z.array(QuesitoEsameDtoSchema).optional()
    },
  }).optional(),
})

export const EsameParzialeParamsDtoSchema = z.object({
  capitoli: z.array(z.number().int()).min(1, 'Seleziona almeno un capitolo'),
  numeroQuesiti: z.number().int().positive(),
  maxErrori: z.number().int().nonnegative(),
  minutiDisponibili: z.number().int().positive(),
}).transform((ui): EsameParzialeParamsDto => ({
  capitoli: ui.capitoli,
  numero_quesiti: ui.numeroQuesiti,
  max_errori: ui.maxErrori,
  minuti_disponibili: ui.minutiDisponibili,
}))

// ==========================================
// 3. Schemi con Trasformazione UI (Remap + CamelCase)
// ==========================================

// --- UI Types (Frontend trasformato) ---

export type User = string

export type Spiegazione = {
  id?: number
  numeroDomanda: number
  spiegazione: string
  focusLinguistico: string
  regolaChiave: string
  domanda?: Domanda
}

export type Argomento = {
  id: number
  nome: string
  domande?: Domanda[]
}

export type Capitolo = {
  id: number
  nome: string
  minNumeroDomanda: number
  maxNumeroDomanda: number
  totaleDomande: number
  domande?: Domanda[]
}

export type CapitoloBasicStats = {
  id: number
  totale: number
  corrette: number
  sbagliate: number
  nonDate: number
  durataMs: number
}

export type Domanda = {
  id: number
  testo: string
  rispostaCorretta: RispostaEnum
  immagine: string | null
  immaginePath: string | null
  idCapitolo: number
  paginaQuiz: number
  idBlocco: number
  argomenti?: Argomento[]
  capitolo?: Capitolo
  spiegazione?: Spiegazione[]
}

export type AttivitaQuesitoEsame = {
  id?: number
  tipo: TipoAttivitaEnum
  rispostaData?: RispostaEnum | null
  inizio: Date
  durataMs: number
  timestamp?: Date
  quesitoEsame?: QuesitoEsame
}

export interface PausaEvent {
  inizio: Date
  fine: Date
}

export type QuesitoEsame = {
  id?: number
  rispostaFinale?: RispostaEnum | null
  createdAt: Date
  updatedAt: Date
  esame?: Esame
  domandaOriginale?: Domanda
  attivita?: AttivitaQuesitoEsame[]
  haSbagliato: boolean
}

export type QuesitiBasicStats = {
  totale: number
  corrette: number
  sbagliate: number
  nonDate: number
}

export type Esame = {
  id?: number
  tipo: TipoEsameEnum
  numeroQuesiti: number
  maxErrori: number
  minutiDisponibili: number
  createdAt: Date
  updatedAt: Date
  utente?: any
  quesiti?: QuesitoEsame[]
  erroriTotali: number
  esitoSuperato: boolean
}

export type EsameParzialeParamsInput = z.input<typeof EsameParzialeParamsDtoSchema>

export const SpiegazioneSchema = SpiegazioneDtoSchema.transform((dto): Spiegazione => ({
  id: dto.id,
  numeroDomanda: dto.numero_domanda,
  spiegazione: dto.spiegazione,
  focusLinguistico: dto.focus_linguistico,
  regolaChiave: dto.regola_chiave,
  domanda: dto.edges?.domanda ? DomandaSchema.parse(dto.edges.domanda) : undefined,
}))

export const ArgomentoSchema = ArgomentoDtoSchema.transform((dto): Argomento => ({
  id: dto.id,
  nome: dto.nome,
  domande: dto.edges?.domande ? dto.edges.domande.map(d => DomandaSchema.parse(d)) : undefined,
}))

export const CapitoloSchema = CapitoloDtoSchema.transform((dto): Capitolo => ({
  id: dto.id,
  nome: dto.nome,
  minNumeroDomanda: dto.min_numero_domanda,
  maxNumeroDomanda: dto.max_numero_domanda,
  totaleDomande: dto.totale_domande,
  domande: dto.edges?.domande ? dto.edges.domande.map(d => DomandaSchema.parse(d)) : undefined,
}))

export const CapitoloBasicStatsSchema = CapitoloBasicStatsDtoSchema.transform(
  (dto): CapitoloBasicStats => ({
    id: dto.id,
    totale: dto.totale,
    corrette: dto.corrette,
    sbagliate: dto.sbagliate,
    nonDate: dto.non_date,
    durataMs: dto.durata_ms,
  }),
)

export const DomandaSchema = DomandaDtoSchema.transform((dto): Domanda => ({
  id: dto.id,
  testo: dto.testo,
  rispostaCorretta: dto.is_true ? RispostaEnum.VERO : RispostaEnum.FALSO,
  immagine: dto.immagine ?? null,
  immaginePath: dto.immagine ? `/quiz_assets/${dto.immagine}.png` : null,
  idCapitolo: dto.id_capitolo,
  paginaQuiz: dto.pagina_quiz,
  idBlocco: dto.id_blocco,
  argomenti: dto.edges?.argomenti ? dto.edges.argomenti.map(a => ArgomentoSchema.parse(a)) : undefined,
  capitolo: dto.edges?.capitolo ? CapitoloSchema.parse(dto.edges.capitolo) : undefined,
  spiegazione: dto.edges?.spiegazione ? dto.edges.spiegazione.map(s => SpiegazioneSchema.parse(s)) : undefined,
}))

export const AttivitaQuesitoEsameSchema = AttivitaQuesitoEsameDtoSchema.transform((dto): AttivitaQuesitoEsame => ({
  id: dto.id,
  tipo: dto.tipo,
  rispostaData: booleanToRisposta(dto.risposta_data),
  inizio: dto.inizio,
  durataMs: dto.durata_ms,
  timestamp: dto.timestamp,
  quesitoEsame: dto.edges?.quesito_esame ? QuesitoEsameSchema.parse(dto.edges.quesito_esame) : undefined,
}))

export const QuesitoEsameSchema = QuesitoEsameDtoSchema.transform((dto): QuesitoEsame => {
  const domanda = dto.edges?.domanda_originale ? DomandaSchema.parse(dto.edges.domanda_originale) : undefined
  const rispostaFinale = booleanToRisposta(dto.risposta_finale)
  let haSbagliato: boolean
  if (!domanda) {
    haSbagliato = false
  } else if (!rispostaFinale) {
    haSbagliato = true
  } else {
    haSbagliato = rispostaFinale !== domanda.rispostaCorretta
  }

  return {
    id: dto.id,
    rispostaFinale,
    createdAt: dto.created_at,
    updatedAt: dto.updated_at,
    esame: dto.edges?.esame ? EsameSchema.parse(dto.edges.esame) : undefined,
    domandaOriginale: domanda,
    attivita: dto.edges?.attivita ? dto.edges.attivita.map(a => AttivitaQuesitoEsameSchema.parse(a)) : undefined,
    haSbagliato,
  }
})

export const QuesitiBasicStatsSchema = QuesitiBasicStatsDtoSchema.transform(
  (dto): QuesitiBasicStats => ({
    totale: dto.totale,
    corrette: dto.corrette,
    sbagliate: dto.sbagliate,
    nonDate: dto.non_date,
  }),
)

export const EsameSchema = EsameDtoSchema.transform((dto): Esame => {
  const quesiti = dto.edges?.quesiti ? dto.edges.quesiti.map(q => QuesitoEsameSchema.parse(q)) : undefined
  const erroriTotali = quesiti ? quesiti.filter(q => q.haSbagliato).length : 0

  return {
    id: dto.id,
    tipo: dto.tipo,
    numeroQuesiti: dto.numero_quesiti,
    maxErrori: dto.max_errori,
    minutiDisponibili: dto.minuti_disponibili,
    createdAt: dto.created_at,
    updatedAt: dto.updated_at,
    utente: dto.edges?.utente,
    quesiti,
    erroriTotali,
    esitoSuperato: erroriTotali <= dto.max_errori,
  }
})



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

const CapitoliSchema = z.array(CapitoloSchema)

export async function getCapitoli (): Promise<Capitolo[]> {
  return hksApi<Capitolo[]>('/capitoli', {
    parseResponse: createParserFor(CapitoliSchema),
  })
}

const CapitoliBasicStatsSchema = z.array(CapitoloBasicStatsSchema)
export async function getCapitoliStats (): Promise<CapitoloBasicStats[]> {
  return hksApi<CapitoloBasicStats[]>(
    '/capitoli/stats',
    {
      parseResponse: createParserFor(CapitoliBasicStatsSchema),
    },
  )
}

export async function spiegaDomandaById (
  domandaId: number,
): Promise<Spiegazione> {
  return hksApi<Spiegazione>(
    `/domande/${domandaId}/spiegazione`,
    { method: 'POST', parseResponse: createParserFor(SpiegazioneSchema) },
  )
}

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

export async function getQuesitiStats (): Promise<QuesitiBasicStats> {
  return hksApi<QuesitiBasicStats>(
    '/esami/quesiti/stats',
    {
      parseResponse: createParserFor(QuesitiBasicStatsSchema),
    },
  )
}

export async function notifyQuesityAttivita (
  quesitoId: number,
  attivita: AttivitaQuesitoEsame,
): Promise<void> {
  let risposta = undefined
  if (attivita.rispostaData) {
    risposta = attivita.rispostaData == RispostaEnum.VERO
  }
  await hksApi(
    `/esami/quesiti/${quesitoId}/attivita`,
    {
      method: 'PUT',
      body: AttivitaQuesitoEsameDtoSchema.parse({
        tipo: attivita.tipo,
        risposta_data: risposta,
        inizio: attivita.inizio.toISOString(),
        durata_ms: attivita.durataMs,
      } as AttivitaQuesitoEsameDto),
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

export async function getEsameById (esameId: number): Promise<Esame> {
  return await hksApi<Esame>(`/esami/${esameId}`, {
    parseResponse: createParserFor(EsameSchema),
  })
}

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
