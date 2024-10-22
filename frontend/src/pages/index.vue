<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

const slugs = ref<string[]>([])
const showModal = ref(false)
const selectedImage = ref<string | null>(null)

const limit = ref(30)
const offset = ref(0)
const loading = ref(false)
const endOfSlugs = ref(false)
const scrollElement = ref<HTMLElement | null>(null)

async function getSlugs() {
  if (loading.value || endOfSlugs.value)
    return

  loading.value = true
  try {
    const response = await fetch(`http://localhost:8080/api/slugs?offset=${offset.value}&limit=${limit.value}`)
    if (response.status === 204) {
      endOfSlugs.value = true
    }
    else {
      const jsonData = await response.json() as string[]
      if (jsonData && jsonData.length > 0) {
        offset.value += jsonData.length
        slugs.value = [...slugs.value, ...jsonData]
      }
      else {
        endOfSlugs.value = true
      }
    }
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
  }
  finally {
    loading.value = false
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

function handleScroll() {
  if (!scrollElement.value)
    return

  const bottomOfElement = scrollElement.value.scrollTop + scrollElement.value.clientHeight >= scrollElement.value.scrollHeight - 100
  if (bottomOfElement) {
    getSlugs()
  }
}

onMounted(() => {
  if (scrollElement.value) {
    scrollElement.value.addEventListener('scroll', handleScroll)
  }
  getSlugs()
})

onUnmounted(() => {
  if (scrollElement.value) {
    scrollElement.value.removeEventListener('scroll', handleScroll)
  }
})
</script>

<template>
  <div ref="scrollElement" class="h-90vh flex flex-col items-center gap-2 overflow-y-auto">
    <div class="flex flex-wrap justify-start gap-2 lg:max-w-8/10">
      <div v-for="(slug, index) in slugs" :key="index">
        <img :src="getThumbnailPath(slug)" :alt="slug" class="h-20vh min-h-200px cursor-pointer" @click="openModal(getImagePath(slug))">
      </div>
    </div>
    <div v-if="loading" class="py-4 text-center">
      Loading...
    </div>
    <div v-if="showModal" class="fixed z-50 size-80% flex items-center justify-center bg-black bg-opacity-75">
      <div class="relative rounded-lg bg-white p-4">
        <button class="absolute right-2 top-2 size-6 rounded-lg bg-red-500 text-white" @click="closeModal">
          X
        </button>
        <img v-if="selectedImage" :src="selectedImage" alt="Selected Image" class="max-h-screen max-w-full">
      </div>
    </div>
  </div>
</template>
