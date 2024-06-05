<template>
  <v-app>
    <v-app-bar app flat>
      <v-spacer></v-spacer>
      <v-btn text="Home" to="/" variant="plain" slim color="primary" />
      <v-btn text="About" to="/about" variant="text" slim/>
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
