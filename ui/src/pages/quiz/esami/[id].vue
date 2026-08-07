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

    <v-card-text>
      <pre>{{ esame }}</pre>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import EsitoIcon from '@/components/EsitoIcon.vue'
  import { type Esame, getEsameById } from '@/services/hankinson'
  import { booleanToChoice } from '@/utils/quesiti'
  import { formatDurationMin, formatDurationMs } from '@/utils/temporal'

  const route = useRoute()

  const esame = ref<Esame>()

  const quesiti = computed(() => {
    return esame.value?.edges?.quesiti
  })

  const numeroPassati = computed(() => {
    return quesiti.value
      ?.filter(quesito => booleanToChoice(quesito.risposta_final) === booleanToChoice(quesito.edges.domanda_originale.is_true))
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
