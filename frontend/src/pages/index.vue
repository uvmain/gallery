<script setup lang="ts">
import { useElementVisibility, useSessionStorage } from '@vueuse/core'
import { computed, onBeforeMount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import { getThumbnailPath } from '../composables/logic'

const router = useRouter()
const slugs = ref<string[]>([])
const showShadowElement = ref()

const selectedImage = useSessionStorage('selected-image', '')

async function getSlugs() {
  try {
    slugs.value = []
    const response = await backendFetchRequest('slugs')
    const jsonData = await response.json() as string[]
    if (jsonData && jsonData.length > 0) {
      slugs.value = [...jsonData]
    }
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
  }
}

function navigateToSlug(slug: string) {
  const slugPath = `/${slug}`
  router.push(slugPath)
}

const headerShadowClass = computed(() => {
  return showShadowElement.value ? ' ' : 'shadow-lg'
})

onBeforeMount(async () => {
  selectedImage.value = undefined
  await getSlugs()
})
</script>

<template>
  <div>
    <Header class="sticky top-0 z-10" :class="headerShadowClass" />
    <div class="flex flex-col items-center p-6">
      <div class="flex flex-wrap gap-x-1 lg:max-w-8/10">
        <div ref="showShadowElement" />
        <div v-for="(slug, index) in slugs" :key="index" class="flex-1 basis-auto">
          <img loading="lazy" :src="getThumbnailPath(slug)" :alt="slug" class="h-full max-h-25vh max-w-40vw min-h-20vh w-full cursor-pointer object-cover" @click="navigateToSlug(slug)">
        </div>
        <div class="flex-2 flex" />
      </div>
    </div>
  </div>
</template>
