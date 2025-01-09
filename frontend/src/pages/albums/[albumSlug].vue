<script setup lang="ts">
import type Dialog from '../../components/Dialog.vue'
import type { Album } from '../../composables/albums'
import { useSessionStorage } from '@vueuse/core'
import { onBeforeMount, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ImageSelector from '../../components/ImageSelector.vue'
import PhotoThumbnail from '../../components/PhotoThumbnail.vue'
import { backendFetchRequest } from '../../composables/fetchFromBackend'
import { getThumbnailPath, niceDate } from '../../composables/logic'

const route = useRoute()
const router = useRouter()

const albumData = ref<Album | undefined>()
const albumLinks = ref<string[]>([])
const deleteDialog = ref<typeof Dialog>()
const addToAlbumDialog = ref<typeof Dialog>()
const albumSlug = ref(route.params.albumSlug as string)
const userLoginState = useSessionStorage('login-state', false)
const selectedImage = useSessionStorage('selected-image', '')
const inEditingMode = ref(false)

const selectedSlugs = useSessionStorage<string[]>('selected-slugs', [])

async function addImagesToAlbum() {
  const newAlbum = {
    AlbumSlug: albumData.value?.Slug,
    ImageSlugs: selectedSlugs.value,
  }
  const options = {
    body: JSON.stringify(newAlbum),
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest('links', options)

  if (response.status === 200) {
    await getLinkData()
  }
}

async function getAlbumData() {
  const response = await backendFetchRequest(`albums/${albumSlug.value}`)
  albumData.value = await response.json()
}

async function getLinkData() {
  const response = await backendFetchRequest(`links/album/${albumSlug.value}`)
  albumLinks.value = await response.json()
}

function showDeleteDialog() {
  deleteDialog.value?.show()
}

async function confirmDeleteAlbum() {
  if (!userLoginState.value) {
    return
  }
  const options = {
    method: 'DELETE',
  }
  try {
    await backendFetchRequest(`albums/${albumSlug.value}`, options)
    router.push('/albums')
  }
  catch (error) {
    console.error('Failed to delete album:', error)
  }
}

function hideDeleteDialog() {
  deleteDialog.value?.hide()
}

function hideAddDialog() {
  addToAlbumDialog.value?.hide()
}

function edit() {
  inEditingMode.value = !inEditingMode.value
}

onBeforeMount(async () => {
  getAlbumData()
  getLinkData()
})
</script>

<template>
  <div class="min-h-screen">
    <Header :show-edit="true" :show-add="true" @edit="edit()" @add="addToAlbumDialog?.show()">
      <div v-if="userLoginState" class="p-2 hover:cursor-pointer" @click="showDeleteDialog">
        <icon-tabler-trash-x class="text-2xl" />
      </div>
    </Header>
    <div v-if="albumData" class="flex flex-row items-center justify-center gap-6 p-6 lg:max-w-8/10">
      <img
        :src="getThumbnailPath(albumData.CoverSlug)"
        :alt="albumData.CoverSlug"
        onerror="this.onerror=null;this.src='/default-image.jpg';"
        class="h-40 w-80 cursor-pointer border-2 border-white border-solid object-cover dark:border-neutral-500"
      />
      <div class="flex flex-col gap-2">
        <div class="text-2xl">
          {{ albumData.Name }}
        </div>
        <div>
          Created: {{ niceDate(albumData.DateCreated) }}
        </div>
        <div v-if="albumLinks">
          {{ albumLinks.length }} photos
        </div>
      </div>
    </div>

    <div id="main" class="grid grid-cols-2 mx-auto gap-8 p-6 lg:grid-cols-7 md:grid-cols-4 lg:max-w-8/10">
      <div v-for="(imageSlug, index) in albumLinks" :key="index" class="relative">
        <PhotoThumbnail :slug="imageSlug" />
      </div>
    </div>

    <Dialog ref="deleteDialog" :close-button="false">
      <div class="m-6 flex flex-col items-center gap-4">
        <icon-tabler-exclamation-circle class="text-4xl text-red" />
        <p class="text-center font-semibold">
          Are you sure you want to delete this album?
        </p>
        <div v-if="albumData?.Name">
          {{ albumData.Name }}
        </div>
      </div>
      <div class="flex justify-center gap-4">
        <button aria-label="cancel" class="button" @click="hideDeleteDialog()">
          Cancel
        </button>
        <button aria-label="delete" class="button" @click="confirmDeleteAlbum()">
          Delete
        </button>
      </div>
    </Dialog>

    <Dialog ref="addToAlbumDialog" :close-button="false" class="size-90%" @keydown.escape="hideAddDialog()">
      <div v-if="albumData" class="flex flex-row items-center justify-center gap-4">
        <img
          :src="getThumbnailPath(albumData.CoverSlug)"
          :alt="albumData.CoverSlug"
          onerror="this.onerror=null;this.src='/default-image.jpg';"
          class="h-40 w-80 cursor-pointer border-2 border-white border-solid object-cover dark:border-neutral-500"
        />
        <div class="text-lg">
          {{ albumData.Name }}
        </div>
        <button aria-label="cancel" class="button" @click="hideAddDialog()">
          Cancel
        </button>
        <button aria-label="delete" class="button" @click="addImagesToAlbum()">
          Add Selected
        </button>
      </div>
      <ImageSelector />
    </Dialog>
  </div>
</template>
