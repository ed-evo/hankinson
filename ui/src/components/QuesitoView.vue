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
              @click="giveAnsware(Choice.VERO)"
              :color="value == Choice.VERO ? 'primary' : undefined"
              :disabled="value == Choice.VERO"
              >V</v-btn
            >
          </v-col>
          <v-col cols="6">
            <v-btn
              icon
              @click="giveAnsware(Choice.FALSO)"
              :color="value == Choice.FALSO ? 'primary' : undefined"
              :disabled="value == Choice.FALSO"
              >F</v-btn
            >
          </v-col>
          <v-col
            cols="12"
            class="align-self-end"
          >
          <slot name="done" :done="done">
              <v-btn
                block
                @click="done"
              >
                Prossimo
              </v-btn>
            </slot>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { Choice, type Domanda, type PausaEvent } from '@/services/hankinson'
import { useDocumentVisibility, useThrottleFn } from '@vueuse/core'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

const props = defineProps<{
  width: number
  height: number
  isLandscape: boolean
  domanda: Domanda
  initialValue?: Choice | null
}>()
const emit = defineEmits<{
  (e: 'ready', at: Date): void
  (e: 'answer', at: Date, value: Choice | null): void
  (e: 'done', at: Date): void
  (e: 'pause', event: PausaEvent): void
}>()

const model = ref<Choice>()
const value = computed(() => model.value ?? props.initialValue)

const giveAnsware = useThrottleFn((choice: Choice) => {
  model.value = choice
  emit('answer', new Date(), choice)
}, 1000)

const done = useThrottleFn(() => {
  emit('done', new Date())
}, 300)

const visibility = useDocumentVisibility()
let pauseStartedAt = new Date()

onMounted(() => {
  emit('ready', new Date())
})

onUnmounted(() => {
  model.value = undefined
})

watch(visibility, (current, old) => {
  if (current === 'hidden') {
    pauseStartedAt = new Date()
  }

  if (current === 'visible' && old === 'hidden') {
    emit('pause', { inizio: pauseStartedAt, fine: new Date() })
  }
})
</script>
