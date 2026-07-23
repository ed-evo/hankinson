<template>
    <v-container :style="{ 'height': `${height - 24}px` }">
        <v-row density="compact" class="h-100" :class="{ 'flex-column': !isLandscape }">
            <v-col class="d-flex align-center justify-center" :class="{ 'h-50': !isLandscape }">
                <img v-if="domanda.immagine" :src="`/quiz_assets/${domanda.immagine}.png`"
                    :style="{ 'max-height': `${height - 56}px` }"></img>
            </v-col>

            <v-col :class="{ 'h-50': !isLandscape }">
                <v-row class="h-50" density="compact">
                    <v-col class="d-flex align-center justify-center text-body-large">
                        {{ domanda.testo }}
                    </v-col>
                </v-row>
                <v-row class="h-45">
                    <template v-if="answer === null">
                        <v-col cols="6" class="text-end">
                            <v-btn icon @click="$emit('answer', Choice.VERO)">V</v-btn>
                        </v-col>
                        <v-col cols="6">
                            <v-btn icon @click="$emit('answer', Choice.FALSO)">F</v-btn>
                        </v-col>
                    </template>
                    <template v-else>
                        <v-col cols="12">
                            Risposta data {{ answer }} è: {{ validateAnsware(domanda, answer) ? 'CORRETTA' :
                            'SBAGLIATA' }}
                        </v-col>
                    </template>
                    <v-col cols="12" class="align-self-end">
                        <v-btn class="w-100" @click="$emit('next')">
                            Prossimo
                        </v-btn>
                    </v-col>
                </v-row>
            </v-col>

        </v-row>

    </v-container>
</template>

<script lang="ts" setup>
import { Choice, type Domanda } from '@/services/hankinson';
import { validateAnsware } from '@/utils/quesiti';

defineProps<{
    width: number,
    height: number,
    isLandscape: boolean,
    domanda: Domanda,
    answer: Choice | null,
}>()
defineEmits<{
    (e: 'answer', choice: Choice | null): void,
    (e: 'next'): void,
}>()
</script>