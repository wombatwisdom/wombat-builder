<template>
  <v-app>
    <v-app-bar app flat>
      <v-spacer></v-spacer>
      <v-btn text="Home" to="/" variant="plain" slim color="primary" />
      <v-btn text="About" to="/about" variant="text" slim/>
      <a
        href="https://github.com/wombatwisdom/wombat-builder"
        title="Wombat Builder"
        class="d-inline-block mx-2 mr-6 social-link"
        rel="noopener noreferrer"
        target="_blank"
      >
        <v-icon
          icon="mdi-github"
        />
      </a>
    </v-app-bar>
    <v-main>
      <router-view />
    </v-main>

    <AppFooter />
  </v-app>
</template>

<script setup>
  //
</script>
