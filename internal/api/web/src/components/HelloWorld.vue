<template>
  <v-container class="fill-height">
    <v-responsive
      class="align-centerfill-height mx-auto"
      max-width="900"
    >
      <div class="my-auto text-center">
        <v-img
          height="250"
          src="@/assets/wombat_logo.svg"
        />

        <h1 class="text-h3 font-weight-bold">wombat</h1>
        <p class="mt-8">
          Wombat is a distribution of <a class="benthos_link" href="https://github.com/redpanda-data/benthos"> RedPanda Benthos</a>, a framework for creating declarative stream processors.
        </p>
      </div>


      <div class="text-center w-50">

      </div>

      <div class="py-8" />

      <div>
        <v-data-table
          hide-default-header
          hide-default-footer
          :items="items">
          <template v-slot:item.os="{ item }">
            <v-icon icon="mdi-apple" v-if="item.os === 'darwin'" />
            <v-icon icon="mdi-linux" v-else-if="item.os === 'linux'" />
            <v-icon icon="mdi-microsoft-windows" v-else-if="item.os === 'windows'" />
          </template>
          <template v-slot:item.link="{ item }">
            <v-btn color="primary" slim download variant="text" :href="`/artifacts/${item.arch}/${item.os}/${goversion}/${pkgHash}`">Download</v-btn>
          </template>
        </v-data-table>
      </div>

    </v-responsive>
  </v-container>
</template>

<script setup>
  import { ref, computed } from 'vue'

  const goversion = ref('1_22_2')
  const pkgHash = ref('e3a62c9055c38740')

  function downloadLink(arch, os) {
    const link = document.createElement('a');
    link.href = `/artifacts/${arch}/${os}/${goversion.value}/${pkgHash.value}`;
    link.target = '_blank';
    link.download = 'my-pdf-file.pdf';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  const items = [
    {
      os: 'darwin',
      arch: 'arm64',
      name: 'Apple with apple silicon',
      link: 'build.arm64.darwin.1_22_2.e3a62c9055c38740',
    },
    {
      os: 'darwin',
      arch: 'amd64',
      name: 'Apple with intel silicon',
      link: 'build.amd64.darwin.1_22_2.e3a62c9055c38740',
    },
    {
      os: 'linux',
      arch: 'arm64',
      name: 'Linux on ARM 64-bit',
      link: 'build.arm64.linux.1_22_2.e3a62c9055c38740',
    },
    {
      os: 'linux',
      arch: 'amd64',
      name: 'Linux on AMD 64-bit',
      link: 'build.amd64.linux.1_22_2.e3a62c9055c38740',
    },
    {
      os: 'windows',
      arch: 'amd64',
      name: 'Windows on AMD 64-bit',
      link: 'build.amd64.windows.1_22_2.e3a62c9055c38740',
    },
  ]
</script>

<style scoped>
  .benthos_link {
    color: red;
    text-decoration: none;
  }
</style>
