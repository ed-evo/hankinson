<template>

  <v-container v-if="domanda" :style="{ 'height': `${height - 24}px` }">
    <v-row density="compact" class="h-100" :class="{ 'flex-column': !isLandscape }">
      <v-col class="d-flex align-center justify-center" :class="{ 'h-50': !isLandscape }">
        <img v-if="domanda.immagine" :src="`/quiz_assets/${domanda.immagine}.png`"
          :style="{ 'max-height': `${height - 56}px` }"></img>
      </v-col>

      <v-col :class="{ 'h-50': !isLandscape }">
        <v-row class="h-50" density="compact">
          <v-col class="d-flex align-center justify-center text-body-large">
            {{ domanda.testo }}
          </v-col>
        </v-row>
        <v-row class="h-45">
          <template v-if="answare === null">
            <v-col cols="6" class="text-end">
              <v-btn icon @click="giveAnsware(Choice.VERO)">V</v-btn>
            </v-col>
            <v-col cols="6">
              <v-btn icon @click="giveAnsware(Choice.FALSO)">F</v-btn>
            </v-col>
          </template>
          <template v-else>
            <v-col cols="12">
              Risposta data {{ answare }} è: {{ validateAnsware(domanda, answare) ? 'CORRETTA' : 'SBAGLIATA' }}
            </v-col>
          </template>
          <v-col cols="12" class="align-self-end">
            <v-btn class="w-100" @click="done">
              Prossimo
            </v-btn>
          </v-col>
        </v-row>
      </v-col>

    </v-row>

  </v-container>
  <v-overlay :model-value="isLoading" class="align-center justify-center" persistent>
    <v-progress-circular indeterminate size="64" />
  </v-overlay>
</template>

<script lang="ts" setup>
import { useQuizStore } from '@/stores/quiz';
import { getDomandaById, nextQuesitoAperto, notifyQuesityAttivita, type Domanda, type Quesito } from '@/services/hankinson';
import { ref, computed, onMounted } from 'vue';

import { useDisplay } from 'vuetify';
import { useThrottleFn } from '@vueuse/core';
import { useAppVisibility } from '@/composables/app';

enum Choice {
  VERO = 'VERO',
  FALSO = 'FALSO'
}

const isLoading = ref(true)
const quizStore = useQuizStore()

const { height, width } = useDisplay()

const { onVisibilityChange } = useAppVisibility()

const isLandscape = computed(() => width.value > height.value)

const quesito = ref<Quesito | null>(null)
const domanda = ref<Domanda | null>(null)

const answare = ref<Choice | null>(null)

const giveAnsware = useThrottleFn((choice: Choice) => {
  if (answare.value == null) {
    answare.value = choice
  }
}, 1000)

function validateAnsware(domanda: Domanda, choice: Choice): boolean {
  console.log(answare, domanda)
  switch (choice) {
    case Choice.VERO:
      return domanda.is_true
    case Choice.FALSO:
      return !domanda.is_true
  }
}

const done = useThrottleFn(async () => {
  isLoading.value = true
  await loadQuesito()
  answare.value = null
  isLoading.value = false
}, 300)

async function loadQuesito() {
  quesito.value = await nextQuesitoAperto(quizStore.capitoliSelezionati)

  domanda.value = await getDomandaById(quesito.value.domandaId)
  
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
    tipo: "pausa",
    inizio: event.startedAt,
    fine: event.finishedAt
  })
})
</script>
