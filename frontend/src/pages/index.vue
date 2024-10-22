<script setup lang="ts">
import { onBeforeMount, ref } from 'vue'

const slugs = ref<string[] | undefined>([])
const showModal = ref(false)
const selectedImage = ref<string | null>(null)

const limit = ref(30)
const offset = ref(0)

async function getSlugs() {
  try {
    fetch(`http://localhost:8080/api/slugs?offset=${offset.value}&limit=${limit.value}`)
      .then(response => response.json())
      .then((json) => {
        const jsonData = json as string[]
        offset.value += jsonData.length
        slugs.value = jsonData
      })
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
  }
}

function openModal(thumbnailPath: string) {
  selectedImage.value = thumbnailPath
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  selectedImage.value = null
}

function getThumbnailPath(slug: string) {
  return `http://localhost:8080/api/thumbnail/${slug}`
}

function getImagePath(slug: string) {
  return `http://localhost:8080/api/optimised/${slug}`
}

onBeforeMount(async () => {
  getSlugs()
})
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-2">
    <div class="flex flex-wrap justify-start gap-2 lg:max-w-8/10">
      <div v-for="(slug, index) in slugs" :key="index">
        <img :src="getThumbnailPath(slug)" :alt="slug" class="h-20vh min-h-200px cursor-pointer" @click="openModal(getImagePath(slug))">
      </div>
    </div>
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75">
      <div class="relative rounded-lg bg-white p-4">
        <button class="absolute right-2 top-2 size-6 rounded-lg bg-red-500 text-white" @click="closeModal">
          X
        </button>
        <img v-if="selectedImage" :src="selectedImage" alt="Selected Image" class="max-h-screen max-w-full">
      </div>
    </div>
  </div>
</template>
