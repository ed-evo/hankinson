<template>
  <v-container :style="{ height: `${height - 24}px` }">
    <v-row
      density="compact"
      class="h-100"
      :class="{ 'flex-column': !isLandscape }"
    >
      <v-col
        class="d-flex align-center justify-center"
        :class="{ 'h-50': !isLandscape }"
      >
        <img
          v-if="domanda.immagine"
          :src="`/quiz_assets/${domanda.immagine}.png`"
          :style="{ 'max-height': `${height - 56}px` }"
        />
      </v-col>

      <v-col :class="{ 'h-50': !isLandscape }">
        <v-row
          class="h-50"
          density="compact"
        >
          <v-col class="d-flex align-center justify-center text-body-large">
            {{ domanda.testo }}
          </v-col>
        </v-row>
        <v-row class="h-45">
          <v-col
            cols="6"
            class="text-end"
          >
            <v-btn
              icon
              @click="answer = Choice.VERO"
              :color="answer == Choice.VERO ? 'primary' : undefined"
              >V</v-btn
            >
          </v-col>
          <v-col cols="6">
            <v-btn
              icon
              @click="answer = Choice.FALSO"
              :color="answer == Choice.FALSO ? 'primary' : undefined"
              >F</v-btn
            >
          </v-col>
          <v-col
            cols="12"
            class="align-self-end"
          >
            <v-btn
              class="w-100"
              @click="$emit('done')"
            >
              Prossimo
            </v-btn>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { Choice, type Domanda, type PausaEvent } from '@/services/hankinson'
import { useDocumentVisibility } from '@vueuse/core'
import { watch } from 'vue'

const answer = defineModel()
defineProps<{
  width: number
  height: number
  isLandscape: boolean
  domanda: Domanda
}>()
const emit = defineEmits<{
  (e: 'done'): void
  (e: 'pause', event: PausaEvent): void
}>()

const visibility = useDocumentVisibility()
let pauseStartedAt = new Date()

watch(visibility, (current, old) => {
  if (current === 'hidden') {
    pauseStartedAt = new Date()
  }

  if (current === 'visible' && old === 'hidden') {
    emit('pause', { inizio: pauseStartedAt, fine: new Date() })
  }
})
</script>
