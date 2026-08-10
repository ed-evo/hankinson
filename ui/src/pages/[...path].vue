<template>
  <v-container class="fill-height">
    <v-responsive class="align-center text-center fill-height">
      <v-empty-state
        headline="404"
        title="Page not found"
        text="The page you are looking for does not exist."
        icon="mdi-alert-circle-outline"
      >
        <template #actions>
          <v-btn
            color="primary"
            variant="flat"
            to="/"
            prepend-icon="mdi-home"
          >
            Back to Home
          </v-btn>
        </template>
      </v-empty-state>
    </v-responsive>
  </v-container>
</template>
