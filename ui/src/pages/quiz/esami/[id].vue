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
  import { ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import { type Esame, getEsameById } from '@/services/hankinson'

  const route = useRoute()

  const esame = ref<Esame>()

  watch(
    () => (route.params as { id: string }).id,
    async newId => {
      esame.value = await getEsameById(Number.parseInt(newId, 10))
      console.log(esame.value)
    },
    { immediate: true },
  )
</script>
