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
        <v-pie
            :items="pieChartData"
            inner-cut="75"
            hide-slice
        >
            <template v-slot:center="{ total }">
                <div class="text-center">
                    <v-chip prepend-icon="mdi-check" color="green"><strong>{{ quesitiStats?.corrette }}</strong></v-chip>
                    <v-chip prepend-icon="mdi-close" color="red">{{ quesitiStats?.sbagliate }}</v-chip>
                </div>
            </template>
        </v-pie>
        {{ quesitiStats }}
    </v-sheet>
</template>

<script lang="ts" setup>
import { VPie } from 'vuetify/labs/VPie'
import { getQuesitiStats, type QuesitiBasicStats } from '@/services/hankinson';
import { useAppStore } from '@/stores/app';
import { useQuizStore } from '@/stores/quiz';
import { computed, onMounted, ref } from 'vue';

const quizStore = useQuizStore()
const appStore = useAppStore()

const quesitiStats = ref<QuesitiBasicStats>()

const pieChartData = computed(() => {
    if (!quesitiStats.value) {
        return undefined
    }
    const {
        corrette,
        sbagliate,
        non_date
    } = quesitiStats.value
    return [
        { key: 'corrette', title: 'Corrette', value: corrette, color: 'green' },
        { key: 'sbagliate', title: 'Sbagliate', value: sbagliate, color: 'red' },
    ]
})

onMounted(async () => {
    quesitiStats.value = await getQuesitiStats()
})
</script>