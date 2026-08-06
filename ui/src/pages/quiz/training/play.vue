<route lang="yaml">
name: quiz_training_play
</route>

<template>
  <v-card
    class="h-100 w-100 position-relative"
    flat
  >
    <v-toolbar ref="cardToolbar">
      <v-btn
        :to="{
          name: 'esame-dettaglio',
          params: { id: quizStore.currentEsameParziale?.id },
        }"
      >Completa</v-btn>

      <v-progress-linear
        absolute
        active
        color="blue"
        height="12"
        location="top"
        :max="tempoMassimo"
        :model-value="timePassed"
        variant="split"
      />
    </v-toolbar>

    <v-window
      v-model="currentIndex"
      :touch="false"
    >
      <v-window-item
        v-for="(current, i) in quizItems"
        :key="i"
      >
        <QuesitoView
          v-if="i === currentIndex"
          :domanda="current.domanda"
          :height="appStore.height - (64 + toolbarHeight)"
          :initial-value="current.answer"
          :is-landscape="appStore.isLandscape"
          :width="appStore.width"
          @answer="(at, choice) => giveAnsware(current, at, choice)"
          @done="(at) => onQuesitoDone(current, at)"
          @pause="onPause(current, $event)"
          @ready="onQuesitoReady"
        >
          <template #done="{ done }">
            <v-row>
              <v-col cols="6">
                <v-btn
                  block
                  @click="changePage(ChangeQuesitoType.PREV, done)"
                >precedente</v-btn>
              </v-col>

              <v-col cols="6">
                <v-btn
                  block
                  @click="changePage(ChangeQuesitoType.NEXT, done)"
                >prossimo</v-btn>
              </v-col>
            </v-row>
          </template>
        </QuesitoView>
      </v-window-item>
    </v-window>
  </v-card>
</template>

<script setup lang="ts">
  import type { Choice, Esame, PausaEvent } from '@/services/hankinson'
  import { useElementSize, useIntervalFn, useThrottleFn } from '@vueuse/core'
  import {
    type ComponentPublicInstance,
    computed,
    onUnmounted,
    ref,
    shallowRef,
    watch,
  } from 'vue'
  import { useRouter } from 'vue-router'
  import QuesitoView from '@/components/QuesitoView.vue'
  import {
    AttivitaEmitter,
    getEsameQuesiti,
    TipoAttivitaQuesito,
  } from '@/services/hankinson'
  import { useAppStore } from '@/stores/app'
  import { useQuizStore } from '@/stores/quiz'
  import { QuizItem } from '@/types/models'
  const router = useRouter()
  const appStore = useAppStore()
  const quizStore = useQuizStore()

  const cardToolbar = ref<ComponentPublicInstance | null>(null)
  const { height: toolbarHeight } = useElementSize(cardToolbar)

  const quizItems = shallowRef<QuizItem[]>([])
  const currentIndex = ref(0)

  const startTime = Date.now()
  const timePassed = ref(0)
  const tempoMassimo = computed(() => {
    const minuti = quizStore.currentEsameParziale?.minuti_disponibili ?? 0
    return minuti * 60 * 1000
  })

  useIntervalFn(() => {
    timePassed.value = Date.now() - startTime
  }, 1000)

  const attivitaEmitter: AttivitaEmitter = new AttivitaEmitter()

  enum ChangeQuesitoType {
    PREV = 'Prev',
    NEXT = 'Next',
    GOTO = 'GoTo',
  }

  const changePage = useThrottleFn(
    (event: ChangeQuesitoType, doneFn?: () => void, newIndex?: number) => {
      let index: number = currentIndex.value
      switch (event) {
        case ChangeQuesitoType.NEXT: {
          index += 1
          break
        }
        case ChangeQuesitoType.PREV: {
          index -= 1
          break
        }
        case ChangeQuesitoType.GOTO: {
          index = newIndex ?? 0
          break
        }
      }
      const targetIndex = Math.min(Math.max(index, 0), quizItems.value.length - 1)
      if (currentIndex.value === targetIndex) {
        return
      }
      currentIndex.value = targetIndex
      if (doneFn) {
        doneFn()
      }
    },
    1000,
  )

  async function onQuesitoReady (at: Date) {
    attivitaEmitter.reset(at)
  }

  const giveAnsware = useThrottleFn(
    async (item: QuizItem, at: Date, choice: Choice | null) => {
      if (item.answer !== choice) {
        item.answer = choice
        await attivitaEmitter.fire(
          item.quesito.id,
          TipoAttivitaQuesito.risposta,
          choice,
          at,
        )
      }
    },
    1000,
  )

  const onQuesitoDone = useThrottleFn(async (item: QuizItem, at: Date) => {
    await attivitaEmitter.fire(
      item.quesito.id,
      item.isAnswered ? TipoAttivitaQuesito.prossimo : TipoAttivitaQuesito.salta,
      null,
      at,
    )
  }, 300)

  async function onPause (item: QuizItem, event: PausaEvent) {
    await attivitaEmitter.firePausa(item.quesito.id, event)
  }

  onUnmounted(() => {
    console.log('unlaoding')
  // quizStore.currentEsameParziale = null
  })

  async function loadQuesiti (esame: Esame) {
    console.log('Fetching quesiti', esame)
    const quesiti = await getEsameQuesiti(esame.id)
    quizItems.value = quesiti.map(quesito => {
      const domanda = quesito.edges.domanda_originale
      return new QuizItem(
        { id: quesito.id, esameId: esame.id, domandaId: domanda.id },
        {
          ...domanda,
        },
      )
    })
  }

  watch(
    () => quizStore.currentEsameParziale,
    async (newEsame, oldEsame) => {
      if (!!newEsame && newEsame !== oldEsame) {
        await loadQuesiti(newEsame)
        return
      }
      router.replace({ name: 'quiz-training' })
    },
    { immediate: true },
  )
</script>
