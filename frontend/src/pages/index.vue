<script setup lang="ts">
import { useIntersectionObserver, useStorage } from '@vueuse/core'
import { computed, onBeforeMount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getServerUrl } from '../composables/getServerBaseUrl'

const router = useRouter()
const observerTarget = ref(null)
const slugs = ref<string[]>([])
const limit = ref(20)
const offset = ref(0)
const loading = ref(false)
const endOfSlugs = ref(false)
const serverBaseUrl = ref()

const state = useStorage('query-state', {
  limit: limit.value,
  offset: offset.value,
  lastSlug: '',
})

async function getSlugs() {
  if (loading.value || endOfSlugs.value)
    return

  loading.value = true
  try {
    const response = await fetch(`${serverBaseUrl.value}/api/slugs?offset=${offset.value}&limit=${limit.value}`)
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

function getThumbnailPath(slug: string) {
  return `${serverBaseUrl.value}/api/thumbnail/${slug}`
}

function navigateToSlug(slug: string) {
  state.value.lastSlug = slug
  const slugPath = `/${slug}`
  router.push(slugPath)
}

const loadingStatus = computed(() => {
  return offset.value > 0 ? 'Loading...' : 'Waiting for initialisation to complete...'
})

onBeforeMount(async () => {
  serverBaseUrl.value = await getServerUrl()
  await getSlugs()
})
</script>

<template>
  <div class="flex flex-col items-center overflow-y-auto bg-gray-100 p-6">
    <div class="flex flex-wrap gap-x-1 lg:max-w-8/10">
      <div v-for="(slug, index) in slugs" :key="index" class="flex-1 basis-auto">
        <img :src="getThumbnailPath(slug)" :alt="slug" class="h-full max-h-25vh max-w-50vw min-h-20vh w-full cursor-pointer object-cover" @click="navigateToSlug(slug)">
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
    <div ref="observerTarget" />
  </div>
</template>
