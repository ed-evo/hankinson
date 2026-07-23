<template>
  <v-app>
    <v-overlay :model-value="quizStore.isLoading" class="align-center justify-center" persistent>
      <div class="text-center">
        <!-- Circular Loading Wheel -->
        <v-progress-circular color="primary" indeterminate size="64" width="6"></v-progress-circular>
        <div class="text-subtitle-1 mt-4 text-white font-weight-medium">
          Loading Quiz from Pi...
        </div>
      </div>
    </v-overlay>
    <v-layout class="d-flex flex-column" v-if="!quizStore.isLoading">
      <v-system-bar color="grey-darken-3">
        <span>
          <v-icon icon="mdi-account"></v-icon>
          {{ quizStore.user }}
        </span>

        <v-spacer></v-spacer>

        <v-icon :icon="settingsDrawerIcon" @click="toggleSettingleDrawer" :color="settingsColor"></v-icon>
      </v-system-bar>
      <v-navigation-drawer v-model="appStore.opennedSettings" color="grey-darken-2" :width="mobile ? width : width / 4">
        <capitoli-select></capitoli-select>
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
import CapitoliSelect from './components/CapitoliSelect.vue';
import { useDisplay } from 'vuetify';
import { useAppStore } from './stores/app.ts';

const { mobile, width } = useDisplay()
const appStore = useAppStore()
const quizStore = useQuizStore()

const settingsDrawerIcon = computed(() =>
  appStore.opennedSettings ? "mdi-close-circle" : "mdi-cog"
)

const settingsColor = computed(() =>
  appStore.opennedSettings ? "red-darken-2" : ""
)

function toggleSettingleDrawer() {
  appStore.opennedSettings = !appStore.opennedSettings
}
</script>
