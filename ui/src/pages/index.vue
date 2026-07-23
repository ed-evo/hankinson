<template>
    <v-sheet color="light-grey" class="h-100 w-100">
        <v-btn :to="{ name: 'quiz-training' }">Start training</v-btn>
        <hr />
        <v-chip
            v-for="capitolo in capitoli" :key="capitolo.id"
        >{{ capitolo.id }}</v-chip>
    </v-sheet>
</template>

<script lang="ts" setup>
import { useQuizStore } from '@/stores/quiz';
import { computed } from 'vue';

const quizStore = useQuizStore()

const capitoli = computed(() => 
    [...quizStore.capitoliSelezionati]
        .sort((a, b) => a - b)
        .map(id => quizStore.capitoli.get(id))
        .filter(capitolo => !!capitolo)
)
</script>