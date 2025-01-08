<script setup lang="ts">
import type Dialog from '../../components/Dialog.vue'
import type { Album } from '../../composables/albums'
import { onBeforeMount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getAllAlbums } from '../../composables/albums'
import { backendFetchRequest } from '../../composables/fetchFromBackend'

const router = useRouter()
const albums = ref()
const addDialog = ref<typeof Dialog>()
const deleteDialog = ref<typeof Dialog>()
const selectedAlbum = ref()
const newAlbumName = ref<string>()

function addAlbum() {
  addDialog.value?.show()
}

async function confirmAddAlbum() {
  const newAlbum = {
    Name: newAlbumName.value,
  }
  const options = {
    body: JSON.stringify(newAlbum),
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest('albums', options)

  hideAddDialog()

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
  hideDeleteDialog()
}

function trashAlbum(album: Album) {
  selectedAlbum.value = album
  deleteDialog.value?.show()
}

function hideAddDialog() {
  addDialog.value?.hide()
}

function hideDeleteDialog() {
  selectedAlbum.value = undefined
  deleteDialog.value?.hide()
}

function navigateToAlbum(albumName: string) {
  router.push(`/albums/${albumName}`)
}

onBeforeMount(async () => {
  albums.value = await getAllAlbums()
})
</script>

<template>
  <div class="min-h-screen">
    <Header :show-add="true" @add="addAlbum" />
    <div id="main" class="grid grid-cols-2 mx-auto gap-8 p-6 lg:grid-cols-7 md:grid-cols-4 lg:max-w-8/10">
      <div v-for="(album, index) in albums" :key="index" class="relative">
        <AlbumCoverLarge :album="album" @trash="trashAlbum(album)" @navigate="navigateToAlbum(album.Slug)" />
      </div>
    </div>
    <Dialog ref="addDialog" :close-button="false" class="border-none shadow-2xl standard">
      <div class="m-6 flex flex-col items-center gap-4">
        <icon-tabler-library-plus class="text-4xl text-green" />
        <div class="flex flex-row items-center gap-2">
          <label for="albumname">Album Name:</label>
          <input id="albumname" v-model="newAlbumName" type="text" name="albumname" autocomplete="albumname">
        </div>
      </div>
      <div class="flex justify-center gap-4">
        <button aria-label="cancel" class="px-4 py-2" @click="hideAddDialog()">
          Cancel
        </button>
        <button aria-label="delete" class="px-4 py-2" @click="confirmAddAlbum()">
          Add
        </button>
      </div>
    </Dialog>
    <Dialog ref="deleteDialog" :close-button="false" class="border-none shadow-2xl standard">
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
        <button aria-label="cancel" class="px-4 py-2" @click="hideDeleteDialog()">
          Cancel
        </button>
        <button aria-label="delete" class="px-4 py-2" @click="deleteAlbum()">
          Delete
        </button>
      </div>
    </Dialog>
  </div>
</template>
