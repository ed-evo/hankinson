<template>
    <v-container>
        <v-btn v-if="!correzioni.length" block :loading="loading" @click="correggi">Correggi</v-btn>
        <correzione-md-view
            v-if="correzioni.length"
            v-for="correzione in correzioni"
            :key="correzione.id"
            :md="correzione.testo"
        ></correzione-md-view>
    </v-container>
</template>
<script setup lang="ts">
import CorrezioneMdView from './CorrezioneMdView.vue';
import { aiCorrege, getCorrezioniEsame, type Correzione } from '@/services/hankinson';
import { useThrottleFn } from '@vueuse/core';
import { onMounted, ref } from 'vue';

const { esameId } = defineProps<{
    esameId: number
}>()

const correzioni = ref<Correzione[]>([])
const loading = ref(false)

const correggi = useThrottleFn(async () => {
    loading.value = true
    try {
        correzioni.value = await aiCorrege(esameId)
    } finally {
        loading.value = false
    }
}, 1000)

async function loadCorrezioniForEsame(id: number) {
    correzioni.value = await getCorrezioniEsame(id)
}

onMounted(async () => {
    await loadCorrezioniForEsame(esameId)
})
</script>