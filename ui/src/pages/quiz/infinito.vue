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
      v-model="current.answer"
      @update:model-value="giveAnsware"
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
            @click="done"
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

let attivitaEmitter: AttivitaEmitter | null = null

const giveAnsware = useThrottleFn((choice: Choice | null) => {
  const item = current.value
  if (item === undefined) {
    return
  }
  if (item.answer !== choice) {
    item.answer = choice
    attivitaEmitter?.fire(TipoAttivitaQuesito.risposta, choice)
  }
}, 1000)

const done = useThrottleFn(async () => {
  isLoading.value = true
  const item = current.value
  attivitaEmitter?.fire(
    item?.answer ? TipoAttivitaQuesito.prossimo : TipoAttivitaQuesito.salta
  )
  await loadQuesito()
  isLoading.value = false
}, 300)

function onPause(event: PausaEvent) {
  console.log('Paused', event)
  attivitaEmitter?.firePausa(event)
}

async function loadQuesito() {
  const quesito = await nextQuesitoAperto(quizStore.capitoliSelezionati)

  const domanda = await getDomandaById(quesito.domandaId)
  const item = new QuizItem(quesito, domanda)
  current.value = item
  quiz.value.push(item)
  attivitaEmitter = new AttivitaEmitter(quesito.id)
}

onMounted(async () => {
  isLoading.value = true
  await loadQuesito()
  isLoading.value = false
})
</script>
