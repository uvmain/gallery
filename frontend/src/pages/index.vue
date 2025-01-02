<script setup lang="ts">
import { useElementVisibility, useSessionStorage } from '@vueuse/core'
import { computed, onBeforeMount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import { getThumbnailPath } from '../composables/logic'

const router = useRouter()
const startObserver = ref<HTMLDivElement | null>(null)
const endObserver = ref<HTMLDivElement | null>(null)
const slugs = ref<string[]>([])
const limit = ref(20)
const offset = ref(0)
const loading = ref(false)
const endOfSlugs = ref(false)

const startObserverIsVisible = useElementVisibility(startObserver)
const endObserverIsVisible = useElementVisibility(endObserver)
const selectedImage = useSessionStorage('selected-image', '')
const state = useSessionStorage('query-state', {
  limit: limit.value,
  offset: offset.value,
  lastSlug: '',
})

async function getSlugs() {
  if (loading.value || endOfSlugs.value)
    return

  loading.value = true
  try {
    const response = await backendFetchRequest(`slugs?offset=${offset.value}&limit=${limit.value}`)
    if (response.status === 204) {
      endOfSlugs.value = true
    }
    else {
      const jsonData = await response.json() as string[]
      if (jsonData && jsonData.length > 0) {
        offset.value += jsonData.length
        slugs.value = [...slugs.value, ...jsonData]

        state.value.limit = limit.value
        state.value.offset = offset.value
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

function navigateToSlug(slug: string) {
  state.value.lastSlug = slug
  const slugPath = `/${slug}`
  router.push(slugPath)
}

const loadingStatus = computed(() => {
  return offset.value > 0 ? 'Loading...' : 'Waiting for initialisation to complete...'
})

watch(endObserverIsVisible, async (newValue) => {
  if (newValue) {
    await getSlugs()
  }
})

const headerShadowClass = computed(() => {
  return startObserverIsVisible.value ? ' ' : 'shadow-lg'
})

onBeforeMount(async () => {
  selectedImage.value = undefined
  await getSlugs()
})
</script>

<template>
  <div>
    <Header class="sticky top-0 z-10" :class="headerShadowClass" />
    <div class="flex flex-col items-center bg-gray-100 p-6">
      <div class="flex flex-wrap gap-x-1 lg:max-w-8/10">
        <div ref="startObserver" />
        <div v-for="(slug, index) in slugs" :key="index" class="flex-1 basis-auto">
          <img :src="getThumbnailPath(slug)" :alt="slug" class="h-full max-h-25vh max-w-40vw min-h-20vh w-full cursor-pointer object-cover" @click="navigateToSlug(slug)">
        </div>
        <div class="flex-2 flex" />
      </div>
      <div v-if="loading" class="py-4 text-center">
        <p>
          {{ loadingStatus }}
        </p>
        <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24">
          <path fill="currentColor" d="M10.72,19.9a8,8,0,0,1-6.5-9.79A7.77,7.77,0,0,1,10.4,4.16a8,8,0,0,1,9.49,6.52A1.54,1.54,0,0,0,21.38,12h.13a1.37,1.37,0,0,0,1.38-1.54,11,11,0,1,0-12.7,12.39A1.54,1.54,0,0,0,12,21.34h0A1.47,1.47,0,0,0,10.72,19.9Z">
            <animateTransform attributeName="transform" dur="0.75s" repeatCount="indefinite" type="rotate" values="0 12 12;360 12 12" />
          </path>
        </svg>
      </div>
      <div ref="endObserver" />
    </div>
  </div>
</template>
