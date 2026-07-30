<route lang="yaml">
name: quiz-infinito
meta:
  type: quiz
  title: Infinito
  icon: mdi-infinity
</route>

<template>
  <v-card class="h-100 w-100 position-relative" flat>

    <v-fab key="absolute" absolute color="primary" location="top left" size="large" icon
      @click="showSpiegazione = !showSpiegazione"
    >
      <v-icon>{{ 'mdi-creation' }}</v-icon>
    </v-fab>
    <v-dialog v-model="showSpiegazione">
      <template v-slot:default="{ isActive }">
        <spiegazione-domanda
          v-if="isActive && quesito?.domandaId"
          :numero-domanda="quesito?.domandaId"
        ></spiegazione-domanda>
      </template>
    </v-dialog>
    <QuesitoView v-if="domanda" :width="appStore.width" :height="appStore.height" :is-landscape="appStore.isLandscape"
      :domanda="domanda" :answer="answer" @answer="giveAnsware" @next="done" @pause="onPause"></QuesitoView>
    <v-overlay :model-value="isLoading" class="align-center justify-center" persistent>
      <v-progress-circular indeterminate size="64" />
    </v-overlay>
  </v-card>
</template>

<script lang="ts" setup>
import SpiegazioneDomanda from '@/components/SpiegazioneDomanda.vue';
import { useQuizStore } from '@/stores/quiz';
import { AttivitaEmitter, Choice, getDomandaById, nextQuesitoAperto, TipoAttivitaQuesito, type Domanda, type PausaEvent, type Quesito } from '@/services/hankinson';
import { ref, onMounted } from 'vue';

import { useThrottleFn } from '@vueuse/core';
import { useAppStore } from '@/stores/app';
import QuesitoView from '@/components/QuesitoView.vue'

const showSpiegazione = ref(false)
const isLoading = ref(true)
const quizStore = useQuizStore()
const appStore = useAppStore()

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

function onPause(event: PausaEvent) {
  console.log("Paused", event)
  attivitaEmitter?.firePausa(event)
}

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
</script>
