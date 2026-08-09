<route lang="yaml">
name: quiz-training
meta:
  type: quiz
  title: Training
  icon: mdi-handball
</route>
<template>
  <v-card>
    <v-card-title> Impostazioni quiz </v-card-title>

    <v-card-text>
      <training-input />
    </v-card-text>

    <v-card-text />

    <v-card-actions>
      <v-btn
        block
        color="primary"
        variant="tonal"
        @click="startParziale"
      >Inizia</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts" setup>
  import type { EsameParzialeParamsInput } from '@/types/hankinson'
  import { useRouter } from 'vue-router'
  import TrainingInput from '@/components/TrainingInput.vue'
  import {
    createEsameParziale,
  } from '@/services/hankinson'
  import { useQuizStore } from '@/stores/quiz'

  const router = useRouter()
  const quizStore = useQuizStore()

  async function startParziale () {
    const settings = quizStore.trainingSettings
    const params: EsameParzialeParamsInput = {
      capitoli: [...quizStore.capitoliSelezionati],
      numeroQuesiti: settings.numeroQuesiti,
      maxErrori: settings.erroriAmmessi,
      minutiDisponibili: settings.numeroQuesiti * settings.secondiPerDomanda,
    }

    const esame = await createEsameParziale(params)
    quizStore.currentEsameParziale = esame
    router.push({ name: 'quiz_training_play' })
  }
</script>
