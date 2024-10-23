<script setup lang="ts">
import { useIntersectionObserver } from '@vueuse/core'
import { computed, onMounted, ref } from 'vue'

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

const loadingStatus = computed(() => {
  return offset.value > 0 ? 'Loading...' : 'Waiting for initialisation to complete...'
})

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
      <p>{{ loadingStatus }}</p>
      <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24"><path fill="currentColor" d="M10.72,19.9a8,8,0,0,1-6.5-9.79A7.77,7.77,0,0,1,10.4,4.16a8,8,0,0,1,9.49,6.52A1.54,1.54,0,0,0,21.38,12h.13a1.37,1.37,0,0,0,1.38-1.54,11,11,0,1,0-12.7,12.39A1.54,1.54,0,0,0,12,21.34h0A1.47,1.47,0,0,0,10.72,19.9Z"><animateTransform attributeName="transform" dur="0.75s" repeatCount="indefinite" type="rotate" values="0 12 12;360 12 12"/></path></svg>
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
