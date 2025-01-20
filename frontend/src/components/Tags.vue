<script setup lang="ts">
import type { Album } from '../composables/albums'
import { onMounted, ref, watch } from 'vue'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import AlbumCoverSmall from './AlbumCoverSmall.vue'

const props = defineProps({
  imageSlug: { type: String, required: true },
  inEditingMode: { type: Boolean, default: false },
})

const tags = ref<string[]>([])
const debug = ref()

async function getTags() {
  tags.value = []
  const response = await backendFetchRequest(`tags/${props.imageSlug}`)
  debug.value = await response.json() || []
  // albumSlugs.forEach(async (albumSlug) => {
  //   const response = await backendFetchRequest(`albums/${albumSlug}`)
  //   const album: Album = await response.json()
  //   imageAlbums.value.push(album)
  // })
}

watch(
  () => props.imageSlug,
  () => {
    getTags()
  },
)

onMounted(() => {
  getTags()
})
</script>

<template>
  <div>
    {{ debug }}
  </div>
</template>
