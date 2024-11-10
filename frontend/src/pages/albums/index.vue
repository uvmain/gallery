<script setup lang="ts">
import type Dialog from '../../components/Dialog.vue'
import type { Album } from '../../composables/albums'
import { onBeforeMount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getAllAlbums } from '../../composables/albums'
import { backendFetchRequest } from '../../composables/fetchFromBackend'

const router = useRouter()
const albums = ref()
const confirmDialog = ref<typeof Dialog>()
const selectedAlbum = ref()

async function addAlbum() {
  const newAlbum = {
    Name: 'Macro',
    CoverSlug: '1729943979792078600',
  }
  const options = {
    body: JSON.stringify(newAlbum),
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest('albums', options)

  if (response.status === 200) {
    albums.value = await getAllAlbums()
  }
}

async function deleteAlbum() {
  const options = {
    method: 'DELETE',
  }
  try {
    await backendFetchRequest(`albums/${selectedAlbum.value?.Slug}`, options)
  }
  catch (error) {
    console.error('Failed to fetch Albums:', error)
  }
  albums.value = await getAllAlbums()
  hideDialog()
}

function trashAlbum(album: Album) {
  selectedAlbum.value = album
  confirmDialog.value?.show()
  console.log(`trashing album ${album.Name}`)
}

function hideDialog() {
  selectedAlbum.value = undefined
  confirmDialog.value?.hide()
}

function navigateToAlbum(albumName: string) {
  router.push(`/albums/${albumName}`)
}

onBeforeMount(async () => {
  albums.value = await getAllAlbums()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <Header :show-add="true" @add="addAlbum" />
    <div id="main" class="grid grid-cols-2 mx-auto gap-8 p-6 lg:grid-cols-7 md:grid-cols-4 lg:max-w-8/10">
      <div v-for="(album, index) in albums" :key="index" class="relative">
        <AlbumCoverLarge :album="album" @trash="trashAlbum(album)" @navigate="navigateToAlbum(album.Slug)" />
      </div>
    </div>
    <Dialog ref="confirmDialog" :close-button="false" class="border-none shadow-2xl">
      <div class="m-6 flex flex-col items-center gap-4">
        <icon-tabler-exclamation-circle class="text-4xl text-red" />
        <p class="text-center font-semibold">
          Are you sure you want to delete this album?
        </p>
        <div v-if="selectedAlbum">
          {{ selectedAlbum.Name }}
        </div>
      </div>
      <div class="flex justify-center gap-4">
        <button class="px-4 py-2" @click="hideDialog()">
          Cancel
        </button>
        <button class="px-4 py-2" @click="deleteAlbum()">
          Delete
        </button>
      </div>
    </Dialog>
  </div>
</template>
