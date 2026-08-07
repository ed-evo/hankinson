<template>
  <v-icon
    v-bind="$attrs"
    :color="isPassed ? passedColor : failedColor"
  >
    {{ isPassed ? passedIcon : failedIcon }}
  </v-icon>
</template>
<script setup lang="ts">
  const {
    isPassed,
    passedColor = 'success',
    failedColor = 'warning',
    passedIcon = 'mdi-check-decagram',
    failedIcon = 'mdi-close-octagon',
  } = defineProps<{
    isPassed: boolean
    passedColor?: string
    failedColor?: string
    passedIcon?: string
    failedIcon?: string
  }>()
</script>
