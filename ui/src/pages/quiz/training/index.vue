<route lang="yaml">
name: quiz-training
meta:
  type: quiz
  title: Training
  icon: mdi-handball
</route>
<template>

  <v-card>
    <v-card-title>
      Impostazioni quiz
    </v-card-title>
    <v-card-text>
      <training-input></training-input>
    </v-card-text>
    <v-card-text>
    </v-card-text>
    <v-card-actions>
      <v-btn
        variant="tonal"
        color="primary"
        block
      >Inizia</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts" setup>
import TrainingInput from '@/components/TrainingInput.vue';

</script>