<route lang="yaml">
meta:
    title: 'Home'
</route>

<template>
    <v-sheet color="light-grey" class="h-100 w-100">
        <v-btn :to="{ name: 'quiz-infinito' }">Start quiz infinito</v-btn>
        <hr />
        Capitoli selezionati: 
        <v-btn icon flat size="sm" @click="appStore.opennedSettings = true">
            <v-icon>mdi-pencil</v-icon>{{ quizStore.capitoliSelezionati?.length }}
        </v-btn>
        <hr />
        {{ quesitiStats }}
    </v-sheet>
</template>

<script lang="ts" setup>
import { getQuesitiStats, type QuesitiBasicStats } from '@/services/hankinson';
import { useAppStore } from '@/stores/app';
import { useQuizStore } from '@/stores/quiz';
import { onMounted, ref } from 'vue';

const quizStore = useQuizStore()
const appStore = useAppStore()

const quesitiStats = ref<QuesitiBasicStats>()

onMounted(async () => {
    quesitiStats.value = await getQuesitiStats()
})
</script>