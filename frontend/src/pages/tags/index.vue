<script setup lang="ts">
import { backendFetchRequest } from '../../composables/fetchFromBackend'

const tags = ref<string[]>([])

async function getAllTags() {
  try {
    const response = await backendFetchRequest('/tags')
    tags.value = await response.json() as string[]
  }
  catch (error) {
    console.error('Failed to fetch tags:', error)
  }
}

onBeforeMount(async () => {
  await getAllTags()
})
</script>

<template>
  <div class="min-h-screen">
    <Header />
    <div class="flex flex-wrap gap-4 p-6">
      <div v-for="(tag, index) in tags" :key="index">
        <Tag :tag="tag" />
      </div>
    </div>
  </div>
</template>
