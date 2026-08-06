<template>
  <v-list
    v-model:selected="quizStore.capitoliSelezionati"
    density="compact"
    select-strategy="leaf"
  >
    <v-list-item
      v-for="[id, capitolo] in quizStore.capitoli"
      :key="id"
      :title="`${id}: ${capitolo.nome}`"
      :value="id"
    >
      <template #prepend="{ isSelected, select }">
        <v-list-item-action start>
          <v-switch
            color="primary"
            hide-details
            inset="material"
            :model-value="isSelected"
            @update:model-value="(value) => select(!!value)"
          />
        </v-list-item-action>
      </template>
    </v-list-item>
  </v-list>
</template>

<script lang="ts" setup>
  import { useQuizStore } from '@/stores/quiz'

  const quizStore = useQuizStore()
</script>
