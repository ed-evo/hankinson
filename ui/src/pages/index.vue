<route lang="yaml">
meta:
    title: 'Home'
</route>

<template>
    <v-sheet color="light-grey" class="h-100 w-100">
        <v-btn
            v-for="quizRoute in appStore.quizRoutes" :key="quizRoute.path"
            :to="quizRoute.path"
            block
        >
        <v-icon v-if="quizRoute.meta.icon">{{quizRoute.meta.icon}}</v-icon>
        <span class="ml-1">{{ quizRoute.meta.title || quizRoute.name }}</span>
        </v-btn>
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
        <v-card>
            <v-card-title>Errori per capitolo</v-card-title>
            <v-sparkline
                :labels="capitoliIds"
                :model-value="erroriPercentuale"
                smooth="4"
                :gradient="['green', 'red']"
                gradient-direction="bottom"
                line-width="1"
                fill
            ></v-sparkline>

            <v-card-title>Tempo medio per capitolo (Globale: {{ format(tempoMedioGlobale || 0) }}sec.)</v-card-title>
            <v-sparkline
                :labels="capitoliIds"
                :model-value="tempiMedi"
                smooth="4"
                :gradient="['lime', 'blue']"
                gradient-direction="bottom"
                line-width="1"
                fill
            ></v-sparkline>
        </v-card>

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

const capitoliIds = computed(() => {
    return capitoliStats.value?.map(capitolo => capitolo.id)
})

const erroriPercentuale = computed(() => {
    return capitoliStats.value?.map(capitolo => (1 - (capitolo.corrette / capitolo.totale) * 100))
})

const tempoMedioGlobale = computed(() => {
    if (!capitoliStats.value) {
        return
    }
    let countQuesiti = 0
    let sumDurate = 0
    for (const capitolo of capitoliStats.value) {
        countQuesiti += capitolo.totale
        sumDurate += capitolo.durata_ms
    }

    return (sumDurate / countQuesiti) / 1000
})

const tempiMedi = computed(() => {
    return capitoliStats.value?.map(capitolo =>
        (capitolo.durata_ms / capitolo.totale) / 1000
    )
})

onMounted(async () => {
    quesitiStats.value = await getQuesitiStats()
    const cStats = await getCapitoliStats()
    capitoliStats.value = cStats.sort((a, b) => a.id - b.id)
})
</script>