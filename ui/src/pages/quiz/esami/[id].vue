<route lang="yaml">
name: esame-dettaglio
</route>

<template>
<v-card>
    <v-card-text>
        <pre>{{ esame }}</pre>
    </v-card-text>
</v-card>
</template>

<script setup lang="ts">
import { getEsameById, type Esame } from '@/services/hankinson';
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';


const route = useRoute()

const esame = ref<Esame>()


watch(
    () => (route.params as { id: string }).id,
    async (newId) => {
        esame.value = await getEsameById(Number.parseInt(newId, 10))
        console.log(esame.value)
    },
    { immediate: true }
)
</script>