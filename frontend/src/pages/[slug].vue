<script setup lang="ts">
import { computed, onBeforeMount, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const slug = ref(route.params.slug as string)
const metadata = ref()

const optimisedPath = computed(() => {
  return `http://localhost:8080/api/optimised/${slug.value}`
})

const fStop = computed(() => {
  const [first, second] = metadata.value ? metadata.value.fStop.split('/').map(Number) : [1, 1]
  return `ƒ/${(first / second).toPrecision(2)}`
})

const focalLength = computed(() => {
  const [first, second] = metadata.value ? metadata.value.focalLength.split('/').map(Number) : [1, 1]
  return `${(first / second).toPrecision(2)} mm`
})

async function getMetadata() {
  try {
    const response = await fetch(`http://localhost:8080/api/metadata/${slug.value}`)
    metadata.value = await response.json() as string[]
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
  }
}

onBeforeMount(() => {
  getMetadata()
})
</script>

<template>
  <div class="flex flex-col p-6 lg:flex-row">
    <img :src="optimisedPath" class="max-h-90vh max-w-70vw" />
    <div>
      <div class="flex items-center p-6 space-x-2">
        <icon-carbon-aperture class="text-2xl text-black" />
        <span>{{ fStop }}</span>
      </div>
      <div class="flex items-center p-6 space-x-2">
        <icon-tabler-eye-table class="text-2xl text-black" />
        <span>{{ focalLength }}</span>
      </div>
    </div>
  </div>
</template>
