<script setup lang="ts">
import type Dialog from '../../components/Dialog.vue'
import type { Album } from '../../composables/albums'
import { useSessionStorage } from '@vueuse/core'
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
const userLoginState = useSessionStorage('login-state', false)
const inEditingMode = ref(false)

function addAlbum() {
  addDialog.value?.show()
}

function edit() {
  inEditingMode.value = !inEditingMode.value
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
    <Header :show-add="userLoginState" :show-edit="!inEditingMode" @edit="edit()" @add="addAlbum">
      <template #2>
        <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="edit()">
          <icon-tabler-edit-off class="text-2xl text-red-700" />
        </div>
      </template>
    </Header>

    <div id="main" class="grid grid-cols-1 mx-auto gap-2 lg:gap-8 p-6 lg:grid-cols-7 md:grid-cols-4 lg:max-w-8/10">
      <div v-for="(album, index) in albums" :key="index" class="relative">
        <AlbumCoverLarge :album="album" :in-edit-mode="inEditingMode" @trash="trashAlbum(album)" @navigate="navigateToAlbum(album.Slug)" />
      </div>
    </div>
    <Dialog ref="addDialog" :close-button="false" class="modal">
      <div class="p-6">
        <div class="flex flex-col items-center gap-6">
          <icon-tabler-library-plus class="text-4xl text-green" />
          <div class="flex flex-row items-center gap-2">
            <label for="albumname">Album Name:</label>
            <input id="albumname" v-model="newAlbumName" type="text" name="albumname" autocomplete="albumname">
          </div>
          <div class="flex justify-center gap-4">
            <button aria-label="cancel" class="button" @click="hideAddDialog()">
              Cancel
            </button>
            <button aria-label="delete" class="button" @click="confirmAddAlbum()">
              Add
            </button>
          </div>
        </div>
      </div>
    </Dialog>
    <Dialog ref="deleteDialog" :close-button="false" class="modal">
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
        <button aria-label="cancel" class="button" @click="hideDeleteDialog()">
          Cancel
        </button>
        <button aria-label="delete" class="button" @click="deleteAlbum()">
          Delete
        </button>
      </div>
    </Dialog>
  </div>
</template>
