<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const slug = ref(route.params.slug as string)

const optimisedPath = computed(() => {
  return `http://localhost:8080/api/optimised/${slug.value}`
})
</script>

<template>
  <div class="p-6">
    <img :src="optimisedPath" class="max-h-90vh max-w-70vw" />
  </div>
</template>
