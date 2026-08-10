<template>
  <v-card v-bind="$attrs">
    <v-card-text>
      <h3>{{ spiegazione?.regolaChiave }}</h3>

      <dl>
        <dt><strong>Spiegazione:</strong></dt>
        <dd>{{ spiegazione?.spiegazione }}</dd>
        <dt><strong>Focus linguistico</strong></dt>
        <dd>{{ spiegazione?.focusLinguistico }}</dd>
      </dl>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
  import type { Spiegazione } from '@/services/hankinson'
  import { ref, watch } from 'vue'
  import {
    spiegaDomandaById,
  } from '@/services/hankinson'

  const spiegazione = ref<Spiegazione>()

  const { numeroDomanda } = defineProps<{
    numeroDomanda: number
  }>()

  watch(
    () => numeroDomanda,
    async (newId, oldId) => {
      if (newId === oldId) {
        return
      }
      if (!newId) {
        spiegazione.value = undefined
        return
      }
      spiegazione.value = await spiegaDomandaById(newId)
    },
    { immediate: true },
  )
</script>
