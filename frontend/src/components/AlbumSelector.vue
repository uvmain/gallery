<script setup lang="ts">
import type { Album } from '../composables/albums'
import { backendFetchRequest } from '../composables/fetchFromBackend'

defineProps({
  inEditingMode: { type: Boolean, default: false },
})

const emits = defineEmits(['cancel', 'confirm'])

const imageAlbums = ref<Album[]>([])
const selectedAlbums = ref<string[]>([])

async function getAllAlbumsList() {
  imageAlbums.value = []
  selectedAlbums.value = []
  const response = await backendFetchRequest('albums')
  imageAlbums.value = await response.json()
}

function toggleSelectAlbum(album: string) {
  const index = selectedAlbums.value.indexOf(album)
  if (index > -1) {
    selectedAlbums.value.splice(index, 1)
  }
  else {
    selectedAlbums.value.push(album)
  }
}

onMounted(() => {
  getAllAlbumsList()
})
</script>

<template>
  <div class="mx-auto w-full flex flex-col gap-4 p-6 lg:w-1/2">
    <div class="flex flex-col items-center justify-center gap-4 lg:flex-row">
      <button aria-label="cancel" class="button" @click="emits('cancel')">
        Cancel
      </button>
      <button aria-label="delete" class="button" @click="emits('confirm', selectedAlbums)">
        Add to selected albums
      </button>
    </div>
    <div class="mt-4 flex flex-wrap gap-x-2 gap-y-1">
      <div v-for="(album, index) in imageAlbums" :key="index" class="relative">
        <div v-show="selectedAlbums.includes(album.Slug)">
          <icon-tabler-square-plus class="absolute z-100 rounded bg-white bg-op-60 text-2xl text-green -right-1 -top-2" />
        </div>
        <AlbumCoverSmall :album="album" :in-editing-mode="inEditingMode" @image-click="toggleSelectAlbum(album.Slug)" />
      </div>
    </div>
  </div>
</template>
