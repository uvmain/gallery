<script setup lang="ts">
import type Dialog from '../../components/Dialog.vue'
import type { Album } from '../../composables/albums'
import { useSessionStorage } from '@vueuse/core'
import { backendFetchRequest } from '../../composables/fetchFromBackend'
import { getThumbnailPath, niceDate } from '../../composables/logic'

const route = useRoute('/albums/[albumSlug]')
const router = useRouter()

const albumData = ref<Album | undefined>()
const albumLinks = ref<string[]>([])
const deleteDialog = ref<typeof Dialog>()
const addToAlbumDialog = ref<typeof Dialog>()
const coverDialog = ref<typeof Dialog>()
const albumSlug = ref(route.params.albumSlug as string)
const userLoginState = useSessionStorage('login-state', false)
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
    selectedSlugs.value = []
    hideAddDialog()
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

async function deleteLink(imageSlug: string) {
  const data = {
    AlbumSlug: albumData.value?.Slug,
    ImageSlug: imageSlug,
  }
  const options = {
    body: JSON.stringify(data),
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest('link', options)

  if (response.status === 200) {
    await getLinkData()
  }
}

async function updateAlbumName() {
  const data = {
    AlbumSlug: albumData.value?.Slug,
    AlbumName: albumData.value?.Name,
  }
  const options = {
    body: JSON.stringify(data),
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest(`albums/name`, options)
  inEditingMode.value = false
  if (response.status === 200) {
    await getAlbumData()
  }
}

async function swapCover(coverSlug: string) {
  const data = {
    AlbumSlug: albumData.value?.Slug,
    CoverSlug: coverSlug,
  }
  const options = {
    body: JSON.stringify(data),
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest(`albums/cover`, options)
  hideCoverDialog()

  if (response.status === 200) {
    await getAlbumData()
  }
}

function showDeleteDialog() {
  deleteDialog.value?.show()
}

function showCoverDialog() {
  selectedSlugs.value = []
  coverDialog.value?.show()
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
  selectedSlugs.value = []
  deleteDialog.value?.hide()
}

function hideCoverDialog() {
  selectedSlugs.value = []
  coverDialog.value?.hide()
}

function hideAddDialog() {
  addToAlbumDialog.value?.hide()
  selectedSlugs.value = []
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
    <Header :show-add="true" :show-edit="!inEditingMode" @edit="edit()" @add="addToAlbumDialog?.show()">
      <div v-if="userLoginState" class="p-2 hover:cursor-pointer" @click="showDeleteDialog">
        <icon-tabler-trash-x class="text-2xl" />
      </div>
      <template #2>
        <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="edit()">
          <icon-tabler-edit-off class="text-2xl text-red-700" />
        </div>
        <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="updateAlbumName">
          <icon-tabler-checkbox class="text-2xl text-green-700" />
        </div>
      </template>
    </Header>

    <div v-if="albumData" class="flex flex-col items-center justify-center gap-6 p-6 lg:max-w-8/10 lg:flex-row">
      <PhotoThumbnail :slug="albumData.CoverSlug" :edit-mode="inEditingMode" :large="true" @edit-image="showCoverDialog()" />
      <div class="flex flex-col gap-2">
        <div v-if="inEditingMode">
          <input id="imageTitle" v-model="albumData.Name" type="text" class="input" @keypress.enter="updateAlbumName">
        </div>
        <div v-else class="text-2xl">
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

    <div class="mx-auto max-w-8/10 flex flex-wrap gap-x-2 gap-y-1">
      <div v-for="(slug, index) in albumLinks" :key="index">
        <PhotoThumbnail :slug="slug" :delete-mode="inEditingMode" :square="true" @delete-image="deleteLink(slug)" />
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
      <div v-if="albumData" class="flex flex-col items-center justify-center gap-4 lg:flex-row">
        <img
          :src="getThumbnailPath(albumData.CoverSlug)"
          :alt="albumData.CoverSlug"
          loading="lazy"
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

    <Dialog ref="coverDialog" :close-button="false" class="size-90%" @keydown.escape="hideAddDialog()">
      <ImageSelector :single-select="true" :defined-slugs="albumLinks" @close-modal="swapCover" />
    </Dialog>
  </div>
</template>
