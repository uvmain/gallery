<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const album = ref(route.params.albumSlug as string)
</script>

<template>
  <div class="min-h-screen bg-gray-300">
    <Header bg="300" :show-edit="true" />
    <div id="main" class="flex flex-row justify-center gap-8 p-6">
      {{ album }}
    </div>
  </div>
</template>
