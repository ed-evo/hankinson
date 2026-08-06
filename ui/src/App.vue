<template>
  <v-app>
    <v-overlay
      class="align-center justify-center"
      :model-value="quizStore.isLoading"
      persistent
    >
      <div class="text-center">
        <!-- Circular Loading Wheel -->
        <v-progress-circular
          color="primary"
          indeterminate
          size="64"
          width="6"
        />

        <div class="text-subtitle-1 mt-4 text-white font-weight-medium">
          Loading Quiz from Pi...
        </div>
      </div>
    </v-overlay>

    <v-layout
      v-if="!quizStore.isLoading"
      class="d-flex flex-column"
    >
      <v-app-bar>
        <template #prepend>
          <v-app-bar-nav-icon
            :color="settingsColor"
            :icon="settingsDrawerIcon"
            @click="toggleSettingleDrawer"
          />
        </template>

        <v-app-bar-title>
          <v-btn to="/"><v-icon>mdi-home</v-icon></v-btn>

          <template v-if="!appStore.mobile">
            <v-btn
              v-for="quizRoute in appStore.quizRoutes"
              :key="quizRoute.path"
              :to="quizRoute.path"
            >
              <v-icon v-if="quizRoute.meta.icon">{{
                quizRoute.meta.icon
              }}</v-icon>

              <span class="ml-1">{{
                quizRoute.meta.title || quizRoute.name
              }}</span>
            </v-btn>
          </template>
        </v-app-bar-title>

        <v-spacer />

        <template #append="">
          <span>
            {{ quizStore.user }}
          </span>
        </template>
      </v-app-bar>

      <v-navigation-drawer
        v-model="appStore.opennedSettings"
        color="grey-darken-2"
        :width="mobile ? width : width / 4"
      >
        <v-btn
          v-for="quizRoute in appStore.quizRoutes"
          :key="quizRoute.path"
          block
          :to="quizRoute.path"
        >
          <v-icon v-if="quizRoute.meta.icon">{{ quizRoute.meta.icon }}</v-icon>
          <span class="ml-1">{{ quizRoute.meta.title || quizRoute.name }}</span>
        </v-btn>

        <capitoli-select />
      </v-navigation-drawer>
      <!--

    <v-navigation-drawer color="grey-darken-1" width="150" permanent>nav 2</v-navigation-drawer>

    <v-app-bar color="grey" height="48" flat>app bar</v-app-bar>

    <v-navigation-drawer color="grey-lighten-1" location="right" width="150" permanent>nav 3</v-navigation-drawer>

    <v-app-bar color="grey-lighten-2" height="48" location="bottom" flat>footer</v-app-bar> -->

      <v-main>
        <router-view />
      </v-main>
    </v-layout>
  </v-app>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useDisplay } from 'vuetify'
  import { useQuizStore } from '@/stores/quiz'
  import CapitoliSelect from './components/CapitoliSelect.vue'
  import { useAppStore } from './stores/app.ts'

  const { mobile, width } = useDisplay()
  const appStore = useAppStore()
  const quizStore = useQuizStore()

  const settingsDrawerIcon = computed(() =>
    appStore.opennedSettings ? 'mdi-close' : 'mdi-menu',
  )

  const settingsColor = computed(() =>
    appStore.opennedSettings ? 'red-darken-2' : '',
  )

  function toggleSettingleDrawer () {
    appStore.opennedSettings = !appStore.opennedSettings
  }
</script>
