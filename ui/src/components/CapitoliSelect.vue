<template>
  <v-list v-model:selected="quizStore.capitoliSelezionati" select-strategy="leaf" density="compact">
    <v-list-item v-for="[id, capitolo] in quizStore.capitoli" :key="id" :title="`${id}: ${capitolo.nome}`"
      :value="id">
      <template v-slot:prepend="{ isSelected, select }">
        <v-list-item-action start>
          <v-switch :model-value="isSelected" @update:model-value="(value) => select(!!value)" color="primary"
            inset="material" hide-details></v-switch>
        </v-list-item-action>
      </template>
    </v-list-item>
  </v-list>
</template>

<script lang="ts" setup>
import { useQuizStore } from '@/stores/quiz';
import { computed } from 'vue';

const quizStore = useQuizStore()

const capitoli = computed(() => Array.from(quizStore.capitoli.values()))

</script>