// Utilities
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useDisplay } from 'vuetify'

export const useAppStore = defineStore('app', () => {
  const router = useRouter()
  const { height, width, mobile } = useDisplay()

  return {
    // states
    opennedSettings: ref(false),

    // getters
    height: computed(() => height.value),
    width: computed(() => width.value),
    mobile: computed(() => mobile.value),
    isLandscape: computed(() => width.value > height.value),
    quizRoutes: computed(() =>
      router?.getRoutes().filter((r) => r.meta.type === 'quiz')
    ),
  }
})
