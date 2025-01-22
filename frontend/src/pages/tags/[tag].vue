<script setup lang="ts">
import { useElementVisibility } from '@vueuse/core'
import { computed, onBeforeMount, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { backendFetchRequest } from '../../composables/fetchFromBackend'
import { getThumbnailPath } from '../../composables/logic'
import { getAllSlugs } from '../../composables/slugs'

const route = useRoute()
const router = useRouter()
const tag = ref(route.params.tag as string)
const startObserver = ref<HTMLDivElement | null>(null)
const slugs = ref<string[]>([])
const startObserverIsVisible = useElementVisibility(startObserver)

async function getSlugs() {
  slugs.value = []
  try {
    const response = await backendFetchRequest(`slugs/tag/${tag.value}`)
    const responseArray = await response.json() as string[]
    if (responseArray != null) {
      slugs.value = slugs.value.concat(responseArray)
    }
  }
  catch (error) {
    console.error('Failed to fetch slugs for tag:', error)
  }
}

function navigateToSlug(slug: string) {
  const slugPath = `/${slug}`
  router.push(slugPath)
}

const headerShadowClass = computed(() => {
  return startObserverIsVisible.value ? ' ' : 'shadow-lg'
})

onBeforeMount(async () => {
  await getSlugs()
  getAllSlugs()
})
</script>

<template>
  <div>
    <Header class="sticky top-0 z-10" :class="headerShadowClass" />
    <div class="flex flex-col items-center p-6">
      <div ref="startObserver" />
      <div class="flex flex-col gap-4 lg:max-w-8/10 lg:flex-row lg:flex-wrap lg:gap-x-2 lg:gap-y-1">
        <div v-for="(slug, index) in slugs" :key="index" class="flex-1 basis-auto">
          <img :src="getThumbnailPath(slug)" :alt="slug" class="h-full min-h-20vh w-full cursor-pointer object-cover lg:max-h-25vh lg:max-w-40vw" @click="navigateToSlug(slug)">
        </div>
        <div class="flex-2 flex" />
      </div>
    </div>
  </div>
</template>
