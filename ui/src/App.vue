<template>
  <v-app>
    <v-overlay
      :model-value="isLoading"
      class="align-center justify-center"
      persistent
    >
      <div class="text-center">
        <!-- Circular Loading Wheel -->
        <v-progress-circular
          color="primary"
          indeterminate
          size="64"
          width="6"
        >{{ quizStore.downloadProgress }}</v-progress-circular>
        <div class="text-subtitle-1 mt-4 text-white font-weight-medium">
          Loading Quiz from Pi...
        </div>
      </div>
    </v-overlay>
    <v-layout class="d-flex flex-column" v-if="!isLoading">
      <v-system-bar color="grey-darken-3">
        <v-icon icon="mdi-wifi-strength-4"></v-icon>
      </v-system-bar>
      <v-navigation-drawer color="grey-darken-2" width="250">
        <v-list lines="one">
          <v-list-item v-for="[id, capitolo ] in quizStore.capitoli.entries()" :key="id"
            :title="`${capitolo.id}: ${capitolo.nome}`"></v-list-item>
        </v-list>
        {{ quizStore.domandeByCapitoli }}
      </v-navigation-drawer>
      <!-- 

    <v-navigation-drawer color="grey-darken-1" width="150" permanent>nav 2</v-navigation-drawer>

    <v-app-bar color="grey" height="48" flat>app bar</v-app-bar>

    <v-navigation-drawer color="grey-lighten-1" location="right" width="150" permanent>nav 3</v-navigation-drawer>

    <v-app-bar color="grey-lighten-2" height="48" location="bottom" flat>footer</v-app-bar> -->

      <v-main class="d-flex flex-column">
        <router-view />
      </v-main>
    </v-layout>
  </v-app>
</template>

<script lang="ts" setup>
import { useQuizStore } from '@/stores/quiz';
import { computed } from 'vue';

const quizStore = useQuizStore()

const isLoading = computed(() => quizStore.downloadProgress < 100)
</script>
