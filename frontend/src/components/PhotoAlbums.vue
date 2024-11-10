<script setup lang="ts">
import type { Album } from '../composables/albums'
import { onMounted, ref, watch } from 'vue'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import AlbumCoverSmall from './AlbumCoverSmall.vue'

const props = defineProps({
  imageSlug: { type: String, required: true },
  inEditingMode: { type: Boolean, default: false },
})

const imageAlbums = ref<Album[]>([])

async function getAllAlbumsList() {
  const response = await backendFetchRequest(`links/image/${props.imageSlug}`)
  const albumSlugs: string[] = await response.json()
  albumSlugs.forEach(async (albumSlug) => {
    const response = await backendFetchRequest(`albums/${albumSlug}`)
    const album: Album = await response.json()
    console.log(`Album: ${JSON.stringify(album)}`)
    imageAlbums.value.push(album)
  })
}

watch(
  () => props.imageSlug,
  () => {
    getAllAlbumsList()
  },
)

onMounted(() => {
  getAllAlbumsList()
})
</script>

<template>
  <div class="mt-4 w-full flex flex-col gap-4 border-1 border-gray-200 border-solid p-4">
    <div class="text-left text-lg">
      This photo is in {{ imageAlbums.length }} albums
    </div>
    <div class="grid grid-cols-2 gap-2 gap-4 lg:grid-cols-4 md:grid-cols-3">
      <div v-for="(album, index) in imageAlbums" :key="index" class="flex flex-col gap-2">
        <AlbumCoverSmall :album="album" />
      </div>
    </div>
  </div>
</template>
