// Utilities
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'

export const useAppStore = defineStore('app', () => {
  const { height, width } = useDisplay()

  return {
    // states
    opennedSettings: ref(true),

    // getters
    height: computed(() => height.value),
    width: computed(() => width.value),
    isLandscape: computed(() => width.value > height.value),
  }
})
