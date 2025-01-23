<script setup lang="ts">
import type { Album } from '../composables/albums'
import { onMounted, ref } from 'vue'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const imageAlbums = ref<Album[]>([])

async function getAllAlbumsList() {
  imageAlbums.value = []
  const response = await backendFetchRequest('albums')
  imageAlbums.value = await response.json()
}

onMounted(() => {
  getAllAlbumsList()
})
</script>

<template>
  <div class="mt-4 flex flex-wrap gap-x-2 gap-y-1">
    <div v-for="(album, index) in imageAlbums" :key="index">
      <AlbumCoverSmall :album="album" />
    </div>
  </div>
</template>
