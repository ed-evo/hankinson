<style lang="css" scoped>
/* Optional minimal Vuetify typography tuning */
.markdown-body {
  :deep(h1),
  :deep(h2),
  :deep(h3) {
    margin-top: 1rem;
    margin-bottom: 0.5rem;
    font-weight: 600;
  }

  :deep(p) {
    margin-bottom: 0.75rem;
  }

  :deep(code) {
    background-color: rgba(0, 0, 0, 0.05);
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
  }
}
</style>

<template>
<div class="markdown-body" v-html="htmlContent"></div>
</template>
<script setup lang="ts">
import { computed } from 'vue';
import { parse } from 'marked'
import DOMPurify from 'dompurify'

const { md } = defineProps<{
    md: string
}>()

const htmlContent = computed(() => {
    if (!md) {
        return ''
    }
    const rawHtml = parse(md) as string
    return DOMPurify.sanitize(rawHtml)
})
</script>