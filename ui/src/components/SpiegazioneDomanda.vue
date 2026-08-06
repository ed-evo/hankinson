<template>
  <v-card v-bind="$attrs">
    <v-card-text>
      <h3>{{ spiegazione?.regola_chiave }}</h3>

      <dl>
        <dt><strong>Spiegazione:</strong></dt>
        <dd>{{ spiegazione?.spiegazione }}</dd>
        <dt><strong>Focus linguistico</strong></dt>
        <dd>{{ spiegazione?.focus_linguistico }}</dd>
      </dl>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
  import { ref, watch } from 'vue'
  import {
    spiegaDomandaById,
    type SpiegazioneDomanda,
  } from '@/services/hankinson'

  const spiegazione = ref<SpiegazioneDomanda>()

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
