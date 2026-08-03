<route lang="yaml">
name: quiz_training_play
</route>

<template>
  <v-card
    class="h-100 w-100 position-relative"
    flat
  >
    <v-window v-model="currentIndex">
      <v-window-item
        v-for="(current, i) in quizItems"
        :key="i"
      >
        <QuesitoView
          v-if="i === currentIndex"
          :width="appStore.width"
          :height="appStore.height"
          :is-landscape="appStore.isLandscape"
          :domanda="current.domanda"
          :answer="current.answer"
          @answer="giveAnsware"
          @next="next"
          @pause="onPause"
        ></QuesitoView>
      </v-window-item>
    </v-window>
  </v-card>
</template>

<script setup lang="ts">
import QuesitoView from '@/components/QuesitoView.vue'
import {
  type Esame,
  getEsameQuesiti,
  type PausaEvent,
  Choice,
} from '@/services/hankinson'
import { useAppStore } from '@/stores/app'
import { useQuizStore } from '@/stores/quiz'
import { QuizItem } from '@/types/models'
import { useThrottleFn } from '@vueuse/core'
import { onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const appStore = useAppStore()
const quizStore = useQuizStore()
const quizItems = ref<QuizItem[]>([])
const currentIndex = ref(0)

const next = useThrottleFn(async () => {
  console.log(next)
}, 300)

const giveAnsware = useThrottleFn((choice: Choice | null) => {
  console.log('answered', choice)
}, 1000)

function onPause(event: PausaEvent) {
  console.log('Paused', event)
}

onUnmounted(() => {
  // quizStore.currentEsameParziale = null
})

async function loadQuesiti(esame: Esame) {
  console.log('Fetching quesiti', esame)
  const quesiti = await getEsameQuesiti(esame.id)
  quizItems.value = quesiti.map((quesito) => {
    const domanda = quesito.edges.domanda_originale
    return new QuizItem(
      { id: quesito.id, esameId: esame.id, domandaId: domanda.id },
      {
        ...domanda,
      }
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
  { immediate: true }
)
</script>
