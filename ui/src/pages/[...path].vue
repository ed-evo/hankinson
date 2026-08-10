<template>
  <v-container class="fill-height">
    <v-responsive class="align-center text-center fill-height">
      <v-empty-state
        headline="404"
        icon="mdi-alert-circle-outline"
        text="The page you are looking for does not exist."
        title="Page not found"
      >
        <template #actions>
          <v-btn
            color="primary"
            prepend-icon="mdi-home"
            to="/"
            variant="flat"
          >
            Back to Home
          </v-btn>
        </template>
      </v-empty-state>
    </v-responsive>
  </v-container>
</template>
