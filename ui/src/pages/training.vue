<route>{"name":"quiz-training"}</route>

<template>
  <quesito-component
    v-if="domanda" 
    :width="appStore.width"
    :height="appStore.height"
    :is-landscape="appStore.isLandscape"
    :domanda="domanda"
    :answer="answer"
    @answer="giveAnsware"
    @next="done"
  ></quesito-component>
  <v-overlay :model-value="isLoading" class="align-center justify-center" persistent>
    <v-progress-circular indeterminate size="64" />
  </v-overlay>
</template>

<script lang="ts" setup>
import { useQuizStore } from '@/stores/quiz';
import { AttivitaEmitter, Choice, getDomandaById, nextQuesitoAperto, notifyQuesityAttivita, TipoAttivitaQuesito, type Domanda, type Quesito } from '@/services/hankinson';
import { ref, onMounted } from 'vue';

import { useThrottleFn } from '@vueuse/core';
import { useAppVisibility } from '@/composables/app';
import { useAppStore } from '@/stores/app';
import QuesitoComponent from '@/components/Quesito.vue'

const isLoading = ref(true)
const quizStore = useQuizStore()
const appStore = useAppStore()

const { onVisibilityChange } = useAppVisibility()

const quesito = ref<Quesito | null>(null)
const domanda = ref<Domanda | null>(null)

const answer = ref<Choice | null>(null)

let attivitaEmitter: AttivitaEmitter | null = null

const giveAnsware = useThrottleFn((choice: Choice | null) => {
  if (answer.value !== choice) {
    answer.value = choice
    attivitaEmitter?.fire(TipoAttivitaQuesito.risposta, choice)
  }
}, 1000)

const done = useThrottleFn(async () => {
  isLoading.value = true
  attivitaEmitter?.fire(
    answer.value ? TipoAttivitaQuesito.prossimo : TipoAttivitaQuesito.salta
  )
  await loadQuesito()
  answer.value = null
  isLoading.value = false
}, 300)

async function loadQuesito() {
  quesito.value = await nextQuesitoAperto(quizStore.capitoliSelezionati)

  domanda.value = await getDomandaById(quesito.value.domandaId)
  attivitaEmitter = new AttivitaEmitter(quesito.value.id)
}

onMounted(async () => {
  isLoading.value = true
  await loadQuesito()
  isLoading.value = false
})

onVisibilityChange(event => {
  if (event.state != 'hidden') {
    return
  }
  if (!quesito.value) {
    return
  }
  console.log("Sleeped for: ", event.finishedAt.getTime() - event.startedAt.getTime())
  notifyQuesityAttivita(quesito.value.id, {
    tipo: TipoAttivitaQuesito.pausa,
    inizio: event.startedAt,
    durata_ms: event.finishedAt.getTime() - event.startedAt.getTime()
  })
})
</script>
