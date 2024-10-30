<script setup lang="ts">
import type Dialog from '../../components/Dialog.vue'
import type { Album } from '../../composables/albums'
import dayjs from 'dayjs'
import { onBeforeMount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { backendFetchRequest, getServerUrl } from '../../composables/fetchFromBackend'

const router = useRouter()
const albums = ref()
const serverBaseUrl = ref()
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
    getAlbums()
  }
}

async function getAlbums() {
  try {
    const response = await backendFetchRequest('albums')
    albums.value = await response.json() as Album
  }
  catch (error) {
    console.error('Failed to fetch Albums:', error)
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
  getAlbums()
  hideDialog()
}

function getImageSource(slug = 'none') {
  return slug === 'none' ? '/default-image.jpg' : `${serverBaseUrl.value}/api/optimised/${slug}`
}

function niceDate(dateString: string) {
  return dayjs(dateString).format('DD/MM/YYYY')
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
  serverBaseUrl.value = await getServerUrl()
  await getAlbums()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <Header :show-add="true" @add="addAlbum" />
    <div id="main" class="grid grid-cols-2 mx-auto gap-8 p-6 lg:grid-cols-7 md:grid-cols-4 lg:max-w-8/10">
      <div v-for="(album, index) in albums" :key="index" class="relative">
        <hr class="mx-auto my-2px h-px max-w-60% border-0 bg-gray-400 opacity-60">
        <hr class="mx-auto my-2px h-px max-w-70% border-0 bg-gray-400 opacity-80">
        <hr class="mx-auto my-2px h-px max-w-80% border-0 bg-gray-400">
        <img :src="getImageSource(album.CoverSlug)" onerror="this.onerror=null;this.src='/default-image.jpg';" class="size-60 border-4 border-white border-solid object-cover" @click="navigateToAlbum(album.Slug)" />
        <div class="absolute right-4 top-3 p-2 hover:cursor-pointer" @click="trashAlbum(album)">
          <icon-tabler-trash-x class="text-xl text-white hover:size-115%" />
        </div>
        <div class="absolute bottom-2 left-2 w-auto flex flex-col gap-2 p-2 text-white">
          <div class="[text-shadow:_0_0px_4px_rgb(0_0_0_/_0.8)] text-lg font-semibold">
            {{ album.Name }}
          </div>
          <div class="[text-shadow:_0_0px_4px_rgb(0_0_0_/_0.8)]">
            Created: {{ niceDate(album.DateCreated) }}
          </div>
        </div>
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
