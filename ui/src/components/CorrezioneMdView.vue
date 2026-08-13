<style lang="scss" scoped>
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

  :deep(ul) {
    padding-left: 1rem;
  }

  :deep(ol) {
    padding-left: 1.5rem;
  }

  :deep(> ul) {
    padding-left: 1rem;
    list-style-type: none;
    > li > ul {
      list-style-type: '🠶 ';   
      li > ul {
        list-style-type: none;
        padding-left: 0;
      }
    }
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

const headerListPtn = /(#+)\s+\d+\.\s+/g

const { md } = defineProps<{
    md: string
}>()

const htmlContent = computed(() => {
    if (!md) {
        return ''
    }
    console.log(md)
    const cleanedMd = md.replaceAll(headerListPtn, '$1 ')
    const rawHtml = parse(cleanedMd) as string
    return DOMPurify.sanitize(rawHtml)
})
</script>