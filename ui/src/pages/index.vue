<template>

  <v-container v-if="currentDomanda" :style="{ 'height': `${height - 24}px` }">
    <v-row density="compact" class="h-100" :class="{ 'flex-column': !isLandscape }">
      <v-col class="d-flex align-center justify-center" :class="{ 'h-50': !isLandscape }">
        <img v-if="currentDomanda.immagine" :src="`/quiz_assets/${currentDomanda.immagine}.png`"
          :style="{ 'max-height': `${height - 56}px` }"></img>
      </v-col>

      <v-col :class="{ 'h-50': !isLandscape }">
        <v-row class="h-50" density="compact">
          <v-col class="d-flex align-center justify-center text-body-large">
            {{ currentDomanda.testo }}
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
              Risposta data {{ answare }} è: {{ validateAnsware(currentDomanda, answare) ? 'CORRETTA' : 'SBAGLIATA' }}
            </v-col>
          </template>
          <v-col cols="12" class="align-self-end">
            <v-btn class="w-100" @click="next">
              Prossimo
            </v-btn>
          </v-col>
        </v-row>
      </v-col>

    </v-row>

  </v-container>

</template>

<script lang="ts" setup>
import { useQuizStore } from '@/stores/quiz';
import type { Domanda } from '@/types/hankinson';
import { ref, computed, onMounted } from 'vue';

import { useDisplay } from 'vuetify';

enum Choice {
  VERO = 'VERO',
  FALSO = 'FALSO'
}

const quizStore = useQuizStore()

const { height, width } = useDisplay()

const isLandscape = computed(() => width.value > height.value)

const currentDomanda = ref<Domanda | null>(null)

const answare = ref<Choice | null>(null)

function giveAnsware(choice: Choice) {
  if (answare.value == null) {
    answare.value = choice
  }
}

function validateAnsware(domanda: Domanda, choice: Choice): boolean {
  console.log(answare, domanda)
  switch (choice) {
    case Choice.VERO:
      return domanda.is_true
    case Choice.FALSO:
      return !domanda.is_true
  }
}

function chooseRandomly<T>(list?: T[]): T {
  if (!list) {
    console.error(list)
    throw new Error("Impossibile selezionare da lista vuota.")
  }
  const randomIndex = Math.floor(Math.random() * list.length);
  return list[randomIndex]
}

function next() {
  answare.value = null
  const capitolo = chooseRandomly(quizStore.capitoliSelezionati)
  console.log("capitolo selezionato: ", capitolo)
  const domande = quizStore.domandeByCapitoli.get(capitolo.id)
  currentDomanda.value = chooseRandomly(domande)
  // currentDomanda.value = quizStore.domandeByCapitoli.get(3)?.find(d => d.id == 19999)
  // currentDomanda.value = quizStore.domandeByCapitoli.get(3)?.find(d => d.id == 19925)
  // currentDomanda.value = quizStore.domandeByCapitoli.get(3)?.find(d => d.id == 19714)
}

onMounted(() => {
  next()
  // setInterval(next, 100)
})
</script>
