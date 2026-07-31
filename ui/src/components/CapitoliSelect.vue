<template>
    <v-btn
      v-for="quizRoute in appStore.quizRoutes" :key="quizRoute.path"
      :to="quizRoute.path"
      block
    >
      <v-icon v-if="quizRoute.meta.icon">{{quizRoute.meta.icon}}</v-icon>
      <span class="ml-1">{{ quizRoute.meta.title || quizRoute.name }}</span>
    </v-btn>
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
import { useAppStore } from '@/stores/app';
import { useQuizStore } from '@/stores/quiz';

const appStore = useAppStore()
const quizStore = useQuizStore()

</script>