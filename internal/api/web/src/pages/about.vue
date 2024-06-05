<template>
  <v-container class="fill-height">
    <v-responsive
      class="align-centerfill-height mx-auto"
      max-width="900"
    >
    <v-row>
      <v-col cols="4">
        <v-img
          class="mt-8 mr-14"
          height="250"
          src="@/assets/wombat_logo.svg"
        />
      </v-col>
      <v-col>
        <v-tabs
          v-model="tab"
          color="primary"
        >
          <v-tab text="About" value="option-1"></v-tab>
          <v-tab text="What's next" value="option-2"></v-tab>
          <v-tab text="Who" value="option-3"></v-tab>
        </v-tabs>

        <v-tabs-window v-model="tab">
          <v-tabs-window-item value="option-1">
            <h1 class="mt-6">Wombat</h1>
            <p class="text-body-1 mt-6">
              Even before the aquisition of Benthos by RedPanda, the goal has always been to make the project bigger than
              a single contributor or company. RedPanda Benthos would provide the core in a similar way as Linux provides
              the kernel. While this core or kernel is the hearth and engine of it all, the ecosystem around it is what
              makes it truly powerful. This is where Wombat comes in.
            </p>
            <p class="text-body-1 mt-6">
              Wombat takes the RedPanda Benthos core and adds a set of plugins, configurations, and tools to make it easier
              to get started with stream processing. The goal is to provide a set of batteries included tools that can be
              used to build a wide range of stream processing applications.
            </p>
            <p class="text-body-1 mt-6">
              The first component within Wombat is the Wombat Builder. The Wombat Builder is the website you are currently
              looking at. What you can't see however, is a backend which provides APIs for you to tailor your wombat to
              your needs by selecting only the packages you need.
            </p>
          </v-tabs-window-item>
          <v-tabs-window-item value="option-2">
            <h1 class="mt-6">What's next</h1>
            <p class="text-body-1 mt-6">
              For the time being, we are just serving prebuilt binaries of the latest version of RedPanda Benthos
              combined with the <span class="monospaced">free</span> components from RedPanda Connect. In the back,
              everything is already in place to build more custom binaries, but we need a bit of time to make this
              available.
            </p>

            <p class="text-body-1 mt-6">
              Some of the other features we are thinking about are:
              <ul class="my-2 ml-6">
                <li>A dedicated page per binary with docs and information how the binary was built</li>
                <li>A wombat-contrib project where plugins can land under an open license (MIT)</li>
              </ul>
              Obviously suggestions are welcome!
            </p>
          </v-tabs-window-item>
          <v-tabs-window-item value="option-3">
            <h1 class="mt-6">Who is behind this?</h1>
            <p class="text-body-1 mt-6">
              Well, I would love to be able to say that this project is the result of a large team of people from the
              community, but the truth is that for the time being it is just me, Daan Gerits. However, I am hoping that
              this project will grow and that more people and companies from the community will join in to make this
              project a success.
            </p>

            <p class="text-body-1 mt-6">
              Reach out through <a href="mailto:daan+wombat@shono.io">email</a> if you want to help out. You
              can also find me in the usual channels (Gophers Slack, RedPanda Slack, RedPanda Discord) as <span class="code">Calmera</span>
            </p>
          </v-tabs-window-item>
        </v-tabs-window>
      </v-col>
    </v-row>
    </v-responsive>
  </v-container>
</template>

<script setup>
  const tab = ref('option-1')
</script>

<style scoped>
  .monospaced {
    font-family: monospace;
  }

  .code {
    font-family: monospace;
    background-color: #f3f4f6;
    padding: 0.2em 0.4em;
    border-radius: 0.2em;
  }
</style>
