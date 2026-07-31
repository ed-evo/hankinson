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
        <v-pie :items="pieChartData" inner-cut="75" hide-slice tooltip>
            <template v-slot:center="{ total }">
                <div class="text-center">
                    <v-chip prepend-icon="mdi-check" color="green">
                        <strong>{{ quesitiStats?.corrette }}</strong>
                    </v-chip>
                    <v-chip prepend-icon="mdi-close" color="red">{{ quesitiStats?.sbagliate }}</v-chip>
                    <v-chip prepend-icon="mdi-help" color="grey">{{ quesitiStats?.non_date }}</v-chip>
                </div>
            </template>
        </v-pie>

        <hr />

        <spiegazione-domanda
            :numero-domanda="18315"
        ></spiegazione-domanda>

        <v-expansion-panels variant="accordion">
            <v-expansion-panel v-for="stat in capitoliStats" :key="stat.id">
                <v-expansion-panel-title>CAP: {{ stat.id }}, TOT: {{ stat.totale }}, Media: {{ format(stat.durata_ms / stat.totale) }}</v-expansion-panel-title>
                <v-expansion-panel-text>
                    <ul>
                        <li><strong>Tempo Totale:</strong> <span>{{ format(stat.durata_ms) }}</span></li>
                        <li><strong>Media risposta:</strong> <span>{{ format(stat.durata_ms / stat.totale) }}</span></li>
                        <li><strong>Totale:</strong> <span>{{ stat.totale }}</span></li>
                        <li><strong>Corrette:</strong> <span>{{ stat.corrette }}</span></li>
                        <li><strong>Sbagliate:</strong> <span>{{ stat.sbagliate }}</span></li>
                        <li><strong>Non date:</strong> <span>{{ stat.non_date }}</span></li>
                    </ul>
                    Tempo medio per
                </v-expansion-panel-text>
            </v-expansion-panel>
        </v-expansion-panels>
    </v-sheet>
</template>

<script lang="ts" setup>
import SpiegazioneDomanda from '@/components/SpiegazioneDomanda.vue';
import { VPie } from 'vuetify/labs/VPie'
import { getCapitoliStats, getQuesitiStats, type CapitoloBasicStats, type QuesitiBasicStats } from '@/services/hankinson';
import { useAppStore } from '@/stores/app';
import { useQuizStore } from '@/stores/quiz';
import { computed, onMounted, ref } from 'vue';
import { addMilliseconds, formatDuration, interval, intervalToDuration } from 'date-fns';

const quizStore = useQuizStore()
const appStore = useAppStore()

const quesitiStats = ref<QuesitiBasicStats>()
const capitoliStats = ref<CapitoloBasicStats[]>()

function format(ms: number) {
    const strat = new Date()
    const end = addMilliseconds(strat, ms)
    const i = interval(strat, end)
    const duration = intervalToDuration(i)
    return formatDuration(duration)
}

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
        { key: 'non_date', title: 'Non date', value: non_date, color: 'grey'}
    ]
})

onMounted(async () => {
    quesitiStats.value = await getQuesitiStats()
    const cStats = await getCapitoliStats()
    capitoliStats.value = cStats.sort((a, b) => a.id - b.id)
})
</script>