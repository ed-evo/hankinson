<template>
  <v-combobox v-model="quizStore.capitoliSelezionati" :items="capitoli" item-value="id" chips closable-chips multiple
    clearable>
    <template v-slot:chip="{ props, item }">
      <v-chip v-bind="props" :text="item.nome" label>
        <template v-slot:prepend>
          <div class="me-1">{{ item.id }}</div>
        </template>
        <template v-slot:close>
          <v-icon icon="$close" size="14"></v-icon>
        </template>
      </v-chip>
    </template>
    <template v-slot:item="{ props, item }">
      <v-list-item v-bind="props" :title="undefined">
        <v-list-item-title>{{ item.id }}: {{ item.nome }}</v-list-item-title>
      </v-list-item>
    </template>
  </v-combobox>
</template>

<script lang="ts" setup>
import { useQuizStore } from '@/stores/quiz';
import { computed, onMounted } from 'vue';

const quizStore = useQuizStore()

const capitoli = computed(() => Array.from(quizStore.capitoli.values()))

onMounted(() => {
  if (quizStore.capitoliSelezionati.length < 1) {
    let count = 0;
    for (const capitolo of quizStore.capitoli.values()) {
      if (capitoli) {
        quizStore.capitoliSelezionati.push(capitolo);
      }
      count += 1;
      if (count >= 3) {
        break;
      }
    }

  }
})
</script>