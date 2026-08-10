<route lang="yaml">
name: esame-dettaglio
</route>

<template>
  <v-card v-if="esame">
    <v-card-title>
      <esito-icon v-if="numeroSbagliati" :is-passed="numeroSbagliati <= esame.maxErrori" />
      {{ numeroPassati }} di {{ esame?.numeroQuesiti }} in {{ formatDurationMs(tempoImpiegato ?? 0) }}
    </v-card-title>

    <v-card-subtitle>
      Errori ammessi: {{ esame.maxErrori }}, Tempo massimo: {{ formatDurationMin(esame.minutiDisponibili) }}
    </v-card-subtitle>

    <v-expansion-panels v-if="quesiti" v-model="openedQuesiti" multiple variant="accordion">
      <v-expansion-panel
        v-for="quesito in quesiti"
        :key="quesito.id"
        :value="quesito.id"
      >
        <v-expansion-panel-title>
          <esito-icon :is-passed="isCorrect(quesito)" />
          {{ quesito.domanda?.testo }}
        </v-expansion-panel-title>

        <v-expansion-panel-text v-if="quesito.domanda">
          <v-img v-if="quesito.domanda.immaginePath" :src="quesito.domanda.immaginePath" width="300" />
          <spiegazione-domanda :numero-domanda="quesito.domanda.id" />
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-card>
</template>

<script setup lang="ts">
  import type { AttivitaQuesitoEsame, Domanda, Esame, QuesitoEsame } from '@/services/hankinson'
  import { computed, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import EsitoIcon from '@/components/EsitoIcon.vue'
  import SpiegazioneDomanda from '@/components/SpiegazioneDomanda.vue'
  import { getEsameById, getEsameQuesiti, getQuestitoDomanda, listQuesitoAttivita } from '@/services/hankinson'
  import { formatDurationMin, formatDurationMs } from '@/utils/temporal'

  const route = useRoute()

  type Quesito = QuesitoEsame & {
    haSbagliato: boolean
    domanda: Domanda
    attivita: AttivitaQuesitoEsame[]
  }

  const esame = ref<Esame>()
  const openedQuesiti = ref<number[]>([])
  const quesiti = ref<Quesito[]>()

  function isCorrect (quesito: Quesito) {
    return !quesito.haSbagliato
  }

  const numeroPassati = computed(() => {
    return quesiti.value
      ?.filter(quesito => !quesito.haSbagliato)
      ?.length
  })

  const numeroSbagliati = computed(() => {
    if (!esame.value || numeroPassati.value === undefined) {
      return undefined
    }
    return esame.value?.numeroQuesiti - numeroPassati.value
  })

  const listaAttivitaRisposta = computed(() => quesiti.value?.flatMap(
    quesito => quesito.attivita,
  )?.filter(attivita => !!attivita),
  )

  const tempoImpiegato = computed(() => {
    return listaAttivitaRisposta.value
      ?.map(attivita => attivita.durataMs)
      ?.reduce((acc, t) => acc + t, 0)
  })

  async function loadEsame (id: number) {
    const e = await getEsameById(id)
    esame.value = e
    const quesitiSbagliatiId: number[] = []

    quesiti.value = []

    for (const q of await getEsameQuesiti(id)) {
      const domanda = await getQuestitoDomanda(q.id)
      const listaAttivita = await listQuesitoAttivita(q.id)
      const quesito: Quesito = {
        ...q,
        haSbagliato: domanda.rispostaCorretta !== q.rispostaFinale,
        domanda: domanda,
        attivita: listaAttivita,
      }
      quesiti.value.push(quesito)
      if (quesito.haSbagliato) {
        quesitiSbagliatiId.push(quesito.id)
      }
    }

    console.info('Quesiti', quesiti.value)

    quesiti.value.sort((a, b) => Number(b.haSbagliato) - Number(a.haSbagliato))
  }

  watch(
    () => (route.params as { id: string }).id,
    async newId => {
      await loadEsame(Number(newId))
    },
    { immediate: true },
  )
</script>
