<template>
<v-card>
    <v-card-title>{{ spiegazione?.regola_chiave }}</v-card-title>
    <v-card-text>
        <dl>
            <dt><strong>Spiegazione:</strong></dt><dd>{{ spiegazione?.spiegazione }}</dd>
            <dt><strong>Focus linguistico</strong></dt><dd>{{ spiegazione?.focus_linguistico }}</dd>
        </dl>
    </v-card-text>
</v-card>
</template>

<script setup lang="ts">
import { spiegaDomandaById, type SpiegazioneDomanda } from '@/services/hankinson';
import { onMounted, onUnmounted, ref } from 'vue';


const spiegazione = ref<SpiegazioneDomanda>()

async function loadSpiegazione(id: number) {
    spiegazione.value = await spiegaDomandaById(id)
}

const { numeroDomanda } = defineProps<{
    numeroDomanda: number
}>()

onMounted(async () => {
    console.log("onMounted")
    await loadSpiegazione(numeroDomanda)
    console.log("Spiegazione loaded")
})
onUnmounted(async () => {
    console.log("unmounted")
})
</script>