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

    <v-expansion-panels v-model="openedQuesiti" v-if="quesiti" variant="accordion" multiple>
      <v-expansion-panel
        v-for="quesito in quesiti"
        :key="quesito.id"
        :value="quesito.id"
      >
        <v-expansion-panel-title>
          <esito-icon :is-passed="isCorrect(quesito)" />
          {{ quesito.domandaOriginale?.testo }}
        </v-expansion-panel-title>

        <v-expansion-panel-text v-if="quesito.domandaOriginale">
          <v-img v-if="quesito.domandaOriginale.immaginePath" :src="quesito.domandaOriginale.immaginePath" width="300" />
          <spiegazione-domanda :numero-domanda="quesito.domandaOriginale.id" />
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-card>
</template>

<script setup lang="ts">
  import type { Esame, QuesitoEsame } from '@/types/hankinson'
  import { computed, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import EsitoIcon from '@/components/EsitoIcon.vue'
  import SpiegazioneDomanda from '@/components/SpiegazioneDomanda.vue'
  import { getEsameById } from '@/services/hankinson'
  import { formatDurationMin, formatDurationMs } from '@/utils/temporal'

  const route = useRoute()

  const esame = ref<Esame>()
  const openedQuesiti = ref<number[]>([])
  const quesiti = computed(() => {
    return esame.value?.quesiti?.
    toSorted((a, b) => Number(b.haSbagliato) - Number(a.haSbagliato))
  })

  function isCorrect (quesito: QuesitoEsame) {
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

  async function loadEsame(id: number) {
    const e = await getEsameById(id)
    esame.value = e
    openedQuesiti.value = e.quesiti?.filter(quesiti => quesiti.haSbagliato)
    ?.map(q => q.id ?? 0).filter(id => !!id) ?? []
  }

  watch(
    () => (route.params as { id: string }).id,
    async newId => {
      await loadEsame(Number(newId))
    },
    { immediate: true },
  )
</script>
