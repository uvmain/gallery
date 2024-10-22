<script setup lang="ts">
import { useIntersectionObserver } from '@vueuse/core'
import { onMounted, ref } from 'vue'

const observerTarget = ref(null)

const slugs = ref<string[]>([])
const showModal = ref(false)
const selectedSlug = ref<string | null>(null)

const limit = ref(20)
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
    useIntersectionObserver(
      observerTarget,
      ([{ isIntersecting }]) => {
        if (isIntersecting && !endOfSlugs.value) {
          getSlugs()
        }
      },
    )
  }
}

function openModal(slug: string) {
  selectedSlug.value = slug
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  selectedSlug.value = null
}

function getThumbnailPath(slug: string) {
  return `http://localhost:8080/api/thumbnail/${slug}`
}

function getOptimisedPath(slug: string) {
  return `http://localhost:8080/api/optimised/${slug}`
}

onMounted(() => {
  getSlugs()
})
</script>

<template>
  <div ref="scrollElement" class="flex flex-col items-center gap-2 overflow-y-auto p-8">
    <div class="flex flex-wrap justify-start gap-2 lg:max-w-8/10">
      <div v-for="(slug, index) in slugs" :key="index">
        <img :src="getThumbnailPath(slug)" :alt="slug" class="h-15vh min-h-100px cursor-pointer" @click="openModal(slug)">
      </div>
    </div>
    <div v-if="loading" class="py-4 text-center">
      Loading...
    </div>
    <div v-if="showModal" class="fixed z-50 max-h-80vh flex items-center justify-center bg-black bg-opacity-75">
      <div class="bg-white p-4">
        <button class="absolute right-2 top-2 size-6 rounded bg-red-500 text-white" @click="closeModal">
          X
        </button>
        <img
          v-if="selectedSlug"
          :src="getOptimisedPath(selectedSlug)"
          alt="Selected Image"
          class="max-h-screen max-w-full"
        >
      </div>
    </div>
    <div ref="observerTarget" />
  </div>
</template>
