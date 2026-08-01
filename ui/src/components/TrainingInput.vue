<template>
    <div>
        <div class="text-body-small">Numero quesiti</div>
        <v-slider v-model="quizStore.trainingSettings.numeroQuesiti" :min="TrainingSettings.minQuesiti"
            :max="TrainingSettings.maxQuesiti" step="1">
            <template v-slot:append="props">
                <v-text-field v-model="quizStore.trainingSettings.numeroQuesiti" density="compact" style="width: 70px"
                    type="number" hide-details single-line></v-text-field>
            </template>
        </v-slider>
    </div>

    <div>
        <div class="text-body-small">Errori ammessi</div>
        <v-slider v-model="quizStore.trainingSettings.erroriAmmessi" :min="TrainingSettings.minErrori"
            :max="TrainingSettings.maxErrori" step="1">
            <template v-slot:append>
                <v-text-field v-model="quizStore.trainingSettings.erroriAmmessi" density="compact" style="width: 70px"
                    type="number" hide-details single-line></v-text-field>
            </template>
        </v-slider>
    </div>

    <div>
        <div class="text-body-small">Secondi per domanda</div>
        <v-slider v-model="quizStore.trainingSettings.secondiPerDomanda" :min="TrainingSettings.minSecondi"
            :max="TrainingSettings.maxSecondi" step="5">
            <template v-slot:append>
                <v-text-field v-model="quizStore.trainingSettings.secondiPerDomanda" density="compact"
                    style="width: 70px" type="number" step="5" hide-details single-line></v-text-field>
            </template>
        </v-slider>
    </div>
    <div>
        <div class="text-body-small">Capitoli selezionati:</div>
        <v-btn icon flat size="sm" @click="appStore.opennedSettings = true">
            <v-icon>mdi-pencil</v-icon> {{ quizStore.capitoliSelezionati?.join(' - ') }}
        </v-btn>
    </div>
</template>

<script setup lang="ts">
import { useAppStore } from '@/stores/app';
import { useQuizStore } from '@/stores/quiz';
import { TrainingSettings } from '@/types/models';

const appStore = useAppStore()
const quizStore = useQuizStore()
</script>