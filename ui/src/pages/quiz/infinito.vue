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
      :domanda="current.domanda"
      :height="appStore.height"
      :initial-value="current.answer"
      :is-landscape="appStore.isLandscape"
      :width="appStore.width"
      @answer="giveAnsware"
      @done="done"
      @pause="onPause"
      @ready="onQuesitoReady"
    />

    <v-bottom-sheet
      v-if="current?.isAnswered"
      inset
      model-value
      persistent
      :scrim="current.isCorrect ? 'success' : 'warning'"
      timeout="-1"
    >
      <v-card>
        <v-card-text class="py-0">
          <h3>
            Risposta data {{ current.answer }} è:
            {{ current.isCorrect ? 'CORRETTA' : 'SBAGLIATA' }}
          </h3>
        </v-card-text>

        <spiegazione-domanda
          class="overflow-y-auto"
          flat
          max-height="35vh"
          :numero-domanda="current.domanda.id"
        />

        <v-card-actions>
          <v-btn
            block
            color="primary"
            variant="tonal"
            @click="done()"
          >
            Prossimo
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-bottom-sheet>

    <v-overlay
      class="align-center justify-center"
      :model-value="isLoading"
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
  import { useThrottleFn } from '@vueuse/core'
  import { onMounted, ref } from 'vue'
  import QuesitoView from '@/components/QuesitoView.vue'
  import SpiegazioneDomanda from '@/components/SpiegazioneDomanda.vue'
  import {
    AttivitaEmitter,
    getDomandaById,
    nextQuesitoAperto,
  } from '@/services/hankinson'
  import { useAppStore } from '@/stores/app'
  import { useQuizStore } from '@/stores/quiz'
  import { type PausaEvent, type RispostaEnum, TipoAttivitaEnum } from '@/types/hankinson'
  import { QuizItem } from '@/types/models'

  const isLoading = ref(true)
  const quizStore = useQuizStore()
  const appStore = useAppStore()

  const quiz = ref<QuizItem[]>([])
  const current = ref<QuizItem>()

  const attivitaEmitter: AttivitaEmitter = new AttivitaEmitter()

  async function onQuesitoReady (at: Date) {
    attivitaEmitter.reset(at)
  }

  const giveAnsware = useThrottleFn(async (at: Date, choice: RispostaEnum | null) => {
    const item = current.value
    if (item === undefined) {
      return
    }
    if (item.answer !== choice) {
      item.answer = choice
      await attivitaEmitter?.fire(
        item.quesitoId,
        TipoAttivitaEnum.RISPOSTA,
        choice,
        at,
      )
    }
  }, 1000)

  const done = useThrottleFn(async (at: Date = new Date()) => {
    isLoading.value = true
    const item = current.value
    if (item) {
      await attivitaEmitter?.fire(
        item.quesitoId,
        item.isAnswered
          ? TipoAttivitaEnum.PROSSIMO
          : TipoAttivitaEnum.SALTA,
        null,
        at,
      )
    }
    await loadQuesito()
    isLoading.value = false
  }, 300)

  async function onPause (pauseEvent: PausaEvent) {
    if (current.value) {
      await attivitaEmitter.firePausa(current.value.quesitoId, pauseEvent)
    }
  }

  async function loadQuesito () {
    current.value = undefined
    const quesito = await nextQuesitoAperto(quizStore.capitoliSelezionati)
    if (!quesito.id || !quesito.domandaOriginale) {
      throw new Error('Domanda originame non presente in quesito')
    }
    const item = new QuizItem(quesito.id, quesito.domandaOriginale)
    current.value = item
    quiz.value.push(item)
  }

  onMounted(async () => {
    isLoading.value = true
    await loadQuesito()
    isLoading.value = false
  })
</script>
