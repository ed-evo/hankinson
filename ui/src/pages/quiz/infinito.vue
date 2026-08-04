<route lang="yaml">
name: quiz-infinito
meta:
  type: quiz
  title: Infinito
  icon: mdi-infinity
</route>

<template>
  <v-card
    class="h-100 w-100 position-relative"
    flat
  >
    <QuesitoView
      v-if="current"
      :width="appStore.width"
      :height="appStore.height"
      :is-landscape="appStore.isLandscape"
      :domanda="current.domanda"
      :initial-value="current.answer"
      @ready="onQuesitoReady"
      @answer="giveAnsware"
      @done="done"
      @pause="onPause"
    ></QuesitoView>

    <v-bottom-sheet
      v-if="current?.isAnswered"
      timeout="-1"
      model-value
      persistent
      inset
      :scrim="current.isCorrect ? 'success' : 'warning'"
    >
      <v-card>
        <v-card-text class="py-0">
          <h3>
            Risposta data {{ current.answer }} è:
            {{ current.isCorrect ? 'CORRETTA' : 'SBAGLIATA' }}
          </h3>
        </v-card-text>
        <spiegazione-domanda
          :numero-domanda="current.domanda.id"
          class="overflow-y-auto"
          max-height="35vh"
          flat
        ></spiegazione-domanda>
        <v-card-actions>
          <v-btn
            color="primary"
            variant="tonal"
            @click="done()"
            block
          >
            Prossimo
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-bottom-sheet>
    <v-overlay
      :model-value="isLoading"
      class="align-center justify-center"
      persistent
    >
      <v-progress-circular
        indeterminate
        size="64"
      />
    </v-overlay>
  </v-card>
</template>

<script lang="ts" setup>
import SpiegazioneDomanda from '@/components/SpiegazioneDomanda.vue'
import { useQuizStore } from '@/stores/quiz'
import {
  AttivitaEmitter,
  Choice,
  getDomandaById,
  nextQuesitoAperto,
  TipoAttivitaQuesito,
  type PausaEvent,
} from '@/services/hankinson'
import { ref, onMounted } from 'vue'

import { useThrottleFn } from '@vueuse/core'
import { useAppStore } from '@/stores/app'
import QuesitoView from '@/components/QuesitoView.vue'
import { QuizItem } from '@/types/models'

const isLoading = ref(true)
const quizStore = useQuizStore()
const appStore = useAppStore()

const quiz = ref<QuizItem[]>([])
const current = ref<QuizItem>()

const attivitaEmitter: AttivitaEmitter = new AttivitaEmitter()

async function onQuesitoReady(at: Date) {
  attivitaEmitter.reset(at)
}

const giveAnsware = useThrottleFn(async (at: Date, choice: Choice | null) => {
  const item = current.value
  if (item === undefined) {
    return
  }
  if (item.answer !== choice) {
    item.answer = choice
    await attivitaEmitter?.fire(
      item.quesito.id,
      TipoAttivitaQuesito.risposta,
      choice,
      at
    )
  }
}, 1000)

const done = useThrottleFn(async (at: Date = new Date()) => {
  isLoading.value = true
  const item = current.value
  if (item) {
    await attivitaEmitter?.fire(
      item.quesito.id,
      item.isAnswered
        ? TipoAttivitaQuesito.prossimo
        : TipoAttivitaQuesito.salta,
      null,
      at
    )
  }
  await loadQuesito()
  isLoading.value = false
}, 300)

async function onPause(pauseEvent: PausaEvent) {
  if (current.value) {
    await attivitaEmitter.firePausa(current.value.quesito.id, pauseEvent)
  }
}

async function loadQuesito() {
  current.value = undefined
  const quesito = await nextQuesitoAperto(quizStore.capitoliSelezionati)

  const domanda = await getDomandaById(quesito.domandaId)
  const item = new QuizItem(quesito, domanda)
  current.value = item
  quiz.value.push(item)
}

onMounted(async () => {
  isLoading.value = true
  await loadQuesito()
  isLoading.value = false
})
</script>
