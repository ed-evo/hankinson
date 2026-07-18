// Utilities
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', {
  state: () => ({
    opennedSettings: ref(true)
  }),
})
