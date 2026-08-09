<route lang="yaml">
name: esame-dettaglio
</route>

<template>
  <v-card v-if="esame">
    <v-card-title>
      <esito-icon v-if="numeroSbagliati" :is-passed="numeroSbagliati <= esame.max_errori" />
      {{ numeroPassati }} di {{ esame?.numero_quesiti }} in {{ formatDurationMs(tempoImpiegato ?? 0) }}
    </v-card-title>

    <v-card-subtitle>
      Errori ammessi: {{ esame.max_errori }}, Tempo massimo: {{ formatDurationMin(esame.minuti_disponibili) }}
    </v-card-subtitle>

    <v-expansion-panels v-if="quesiti" variant="accordion">
      <v-expansion-panel
        v-for="quesito in quesiti"
        :key="quesito.id"
      >
        <v-expansion-panel-title>
          <esito-icon :is-passed="isCorrect(quesito)" />
          {{ quesito.edges?.domanda_originale?.testo }}
        </v-expansion-panel-title>

        <v-expansion-panel-text v-if="quesito.edges?.domanda_originale">
          <v-img v-if="quesito.edges.domanda_originale.immagine" :src="getImmaginePath(quesito.edges.domanda_originale)" width="300" />
          <spiegazione-domanda :numero-domanda="quesito.edges.domanda_originale.id" />
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-card>
</template>

<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import EsitoIcon from '@/components/EsitoIcon.vue'
  import SpiegazioneDomanda from '@/components/SpiegazioneDomanda.vue'
  import { type Esame, getImmaginePath, getEsameById, type QuesitoEsame } from '@/services/hankinson'
  import { booleanToChoice } from '@/utils/quesiti'
  import { formatDurationMin, formatDurationMs } from '@/utils/temporal'

  const route = useRoute()

  const esame = ref<Esame>()

  const quesiti = computed(() => {
    return esame.value?.edges?.quesiti
  })

  function isCorrect (quesito: QuesitoEsame) {
    return booleanToChoice(quesito.risposta_finale) === booleanToChoice(quesito.edges?.domanda_originale?.is_true ?? false)
  }

  const numeroPassati = computed(() => {
    return quesiti.value
      ?.filter(quesito => isCorrect(quesito))
      ?.length
  })

  const numeroSbagliati = computed(() => {
    if (!esame.value || numeroPassati.value === undefined) {
      return undefined
    }
    return esame.value?.numero_quesiti - numeroPassati.value
  })

  const listaAttivitaRisposta = computed(() => quesiti.value?.flatMap(
    quesito => quesito.edges?.attivita,
  )?.filter(attivita => !!attivita),
  )

  const tempoImpiegato = computed(() => {
    return listaAttivitaRisposta.value
      ?.map(attivita => attivita.durata_ms)
      ?.reduce((acc, t) => acc + t, 0)
  })

  watch(
    () => (route.params as { id: string }).id,
    async newId => {
      esame.value = await getEsameById(Number.parseInt(newId, 10))
      console.log(esame.value)
    },
    { immediate: true },
  )
</script>
