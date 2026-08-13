<template>
  <v-container>
    <v-btn
      v-if="correzioni.length === 0"
      block
      color="primary"
      :loading="loading"
      variant="tonal"
      @click="correggi"
    >
      <template #loader><v-progress-circular indeterminate /> Può richiedere fino a 2 minuti</template>
      <template #prepend><v-icon>mdi-head-snowflake-outline</v-icon></template>
      Correzione IA
    </v-btn>

    <template v-if="correzioni.length > 0">
      <v-tabs v-model="currentTab">
        <v-tab
          v-for="correzione in correzioni"
          :key="correzione.id"
          :value="correzione.id"
        >{{ correzione.type }}:{{ correzione.esaminatore }}</v-tab>
      </v-tabs>

      <v-tabs-window v-model="currentTab">
        <v-tabs-window-item
          v-for="correzione in correzioni"
          :key="correzione.id"
          :value="correzione.id"
        >
          <correzione-md-view
            :md="correzione.testo"
          />
        </v-tabs-window-item>
      </v-tabs-window>

    </template>
  </v-container>
</template>
<script setup lang="ts">
  import { useThrottleFn } from '@vueuse/core'
  import { onMounted, ref } from 'vue'
  import { aiCorrege, type Correzione, getCorrezioniEsame } from '@/services/hankinson'
  import CorrezioneMdView from './CorrezioneMdView.vue'

  const { esameId } = defineProps<{
    esameId: number
  }>()

  const correzioni = ref<Correzione[]>([])
  const loading = ref(false)
  const currentTab = ref<number>()

  const correggi = useThrottleFn(async () => {
    loading.value = true
    try {
      correzioni.value = await aiCorrege(esameId)
    } finally {
      loading.value = false
    }
  }, 1000)

  async function loadCorrezioniForEsame (id: number) {
    correzioni.value = await getCorrezioniEsame(id)
  }

  onMounted(async () => {
    await loadCorrezioniForEsame(esameId)
  })
</script>
