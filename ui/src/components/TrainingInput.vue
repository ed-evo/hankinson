<template>
  <div>
    <div class="text-body-small">Numero quesiti</div>

    <v-slider
      v-model="quizStore.trainingSettings.numeroQuesiti"
      :max="TrainingSettings.maxQuesiti"
      :min="TrainingSettings.minQuesiti"
      step="1"
    >
      <template #append>
        <v-text-field
          v-model="quizStore.trainingSettings.numeroQuesiti"
          density="compact"
          hide-details
          single-line
          style="width: 70px"
          type="number"
        />
      </template>
    </v-slider>
  </div>

  <div>
    <div class="text-body-small">Errori ammessi</div>

    <v-slider
      v-model="quizStore.trainingSettings.erroriAmmessi"
      :max="TrainingSettings.maxErrori"
      :min="TrainingSettings.minErrori"
      step="1"
    >
      <template #append>
        <v-text-field
          v-model="quizStore.trainingSettings.erroriAmmessi"
          density="compact"
          hide-details
          single-line
          style="width: 70px"
          type="number"
        />
      </template>
    </v-slider>
  </div>

  <div>
    <div class="text-body-small">Secondi per domanda</div>

    <v-slider
      v-model="quizStore.trainingSettings.secondiPerDomanda"
      :max="TrainingSettings.maxSecondi"
      :min="TrainingSettings.minSecondi"
      step="5"
    >
      <template #append>
        <v-text-field
          v-model="quizStore.trainingSettings.secondiPerDomanda"
          density="compact"
          hide-details
          single-line
          step="5"
          style="width: 70px"
          type="number"
        />
      </template>
    </v-slider>
  </div>

  <div>
    <div class="text-body-small">Capitoli selezionati:</div>

    <v-btn
      flat
      icon
      size="sm"
      @click="appStore.opennedSettings = true"
    >
      <v-icon>mdi-pencil</v-icon>
      {{ quizStore.capitoliSelezionati?.join(' - ') }}
    </v-btn>
  </div>
</template>

<script setup lang="ts">
  import { useAppStore } from '@/stores/app'
  import { useQuizStore } from '@/stores/quiz'
  import { TrainingSettings } from '@/types/models'

  const appStore = useAppStore()
  const quizStore = useQuizStore()
</script>
