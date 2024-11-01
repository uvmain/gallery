<script setup lang="ts">
import type { Album } from '../composables/albums'
import type { ImageMetadata } from '../composables/imageMetadataInterface'
import { useStorage } from '@vueuse/core'
import dayjs from 'dayjs'
import { computed, onBeforeMount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getAlbumCoverSlugThumbnailAddress, getAlbums } from '../composables/albums'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const userLoginState = useStorage('login-state', false)

const route = useRoute()

const slug = ref(route.params.imageSlug as string)
const imageSize = ref<string>('optimised')
const imageAlbums = ref<Album[]>([])
const loadOriginalText = ref()
const metadata = ref<ImageMetadata | undefined>()
const inEditingMode = ref(false)

const imageSource = computed(() => {
  return `/api/${imageSize.value}/${slug.value}`
})

const fStop = computed(() => {
  const [first, second] = metadata.value ? metadata.value.fStop.split('/').map(Number) : [1, 1]
  const result = (first / second)
  return Number.isNaN(result) ? undefined : `ƒ/${result.toFixed(1)}`
})

const focalLength = computed(() => {
  const [first, second] = metadata.value ? metadata.value.focalLength.split('/').map(Number) : [1, 1]
  const result = (first / second)
  return Number.isNaN(result) ? undefined : `${result.toFixed(1)} mm`
})

const dateTaken = computed(() => {
  return metadata.value && metadata.value.dateTaken !== 'unknown' ? dayjs(metadata.value.dateTaken).format('DD MMMM YYYY HH:mm:ss') : undefined
})

const dateUploaded = computed(() => {
  return metadata.value && metadata.value.dateUploaded !== 'unknown' ? dayjs(metadata.value.dateUploaded).format('DD MMMM YYYY HH:mm:ss') : undefined
})

const camera = computed(() => {
  return metadata.value && metadata.value.cameraMake !== 'unknown' && metadata.value.cameraModel !== 'unknown' ? `${metadata.value.cameraMake} ${metadata.value.cameraModel}` : undefined
})

const lens = computed(() => {
  if (!metadata.value || metadata.value.lensMake === 'unknown' || metadata.value.lensModel === 'unknown') {
    return undefined
  }
  else if (metadata.value?.lensModel?.split(' ')[0].toLocaleLowerCase() === metadata.value?.lensMake?.split(' ')[0].toLocaleLowerCase()) {
    return metadata.value.lensModel
  }
  else {
    return `${metadata.value.lensMake} ${metadata.value.lensModel}`
  }
})

async function getAlbumsList() {
  const albums = await getAlbums()
  imageAlbums.value = albums
}

const whiteBalance = computed(() => {
  let result
  const setting = metadata.value?.whiteBalance
  const temp = metadata.value?.whiteBalanceMode
  if (setting && setting.toLowerCase() !== 'unknown') {
    result = setting
  }
  if (temp && temp.toLowerCase() !== 'unknown') {
    result = `${setting}: ${temp}`
  }
  return result
})

const loadOriginalIconColour = computed(() => {
  return imageSize.value === 'optimised' ? 'text-gray-600' : 'text-gray-400'
})

async function getMetadata() {
  try {
    const response = await backendFetchRequest(`metadata/${slug.value}`)
    metadata.value = await response.json() as ImageMetadata
    imageSize.value = 'optimised'
    loadOriginalText.value = 'Load Original'
    getAlbumsList()
  }
  catch (error) {
    console.error('Failed to fetch Metadata:', error)
  }
}

function loadOriginal() {
  imageSize.value = 'original'
  loadOriginalText.value = 'Original image loaded'
}

async function downloadOriginal() {
  const response = await backendFetchRequest(`original/${slug.value}`)
  const imageBlob = await response.blob()
  const a = document.createElement('a')
  const url = window.URL.createObjectURL(imageBlob)
  const fileName = metadata.value?.fileName || `${slug.value}.jpeg`
  a.href = url
  a.download = fileName
  a.click()
  window.URL.revokeObjectURL(url)
}

function enableEditing() {
  if (userLoginState.value) {
    inEditingMode.value = true
  }
  else {
    console.warn('Not authorised to enter editing mode, please log in')
  }
}

function disableEditing() {
  inEditingMode.value = false
  getMetadata()
}

async function saveMetadata() {
  if (userLoginState.value) {
    inEditingMode.value = false
    try {
      const options = {
        body: JSON.stringify(metadata.value),
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
      }
      await backendFetchRequest(`metadata/${slug.value}`, options)
      await getMetadata()
    }
    catch (error) {
      console.error('Failed to update Metadata:', error)
    }
  }
  else {
    console.warn('Not authorised to save image metadata, please log in')
  }
}

function addToAlbum() {
  console.log('adding to album')
}

watch(
  () => route.params.imageSlug,
  () => {
    slug.value = route.params.imageSlug as string
    getMetadata()
  },
)

onBeforeMount(async () => {
  await getMetadata()
})
</script>

<template>
  <div class="min-h-screen bg-gray-300">
    <Header bg="300" :show-edit="!inEditingMode" @edit="enableEditing">
      <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="disableEditing">
        <icon-tabler-edit-off class="text-2xl text-gray-700" />
      </div>
      <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="saveMetadata">
        <icon-tabler-checkbox class="text-2xl text-gray-700" />
      </div>
    </Header>
    <hr class="mx-auto mt-2 h-px border-0 bg-gray-400 opacity-60 lg:max-w-7/10">
    <div id="main" class="flex flex-row justify-center gap-8 p-6">
      <!-- Image Section -->
      <img v-if="imageSource" :src="imageSource" class="max-h-80vh max-w-70vw border-6 border-white border-solid" />

      <!-- EXIF Data Section -->
      <div v-if="metadata" class="flex flex-col gap-6 p-6 text-sm lg:max-w-1/3">
        <div>
          <div v-if="inEditingMode">
            <label for="imageTitle" class="mb-4 text-2xl">
              Title:
            </label>
            <input id="imageTitle" v-model="metadata.title" type="text">
          </div>
          <h1 v-else class="mb-4 text-2xl text-gray-700 font-semibold">
            {{ metadata.title }}
          </h1>
        </div>

        <div v-if="camera" class="flex items-center space-x-3">
          <icon-tabler-camera class="text-3xl text-gray-600" />
          <div class="flex flex-col gap-1 text-base">
            <span class="text-gray-700 font-medium">{{ camera }}</span>
            <span class="text-gray-600">{{ lens }}</span>
          </div>
        </div>

        <div>
          <div v-if="inEditingMode">
            <label for="dateTaken" class="mb-4 text-2xl">
              Date taken:
            </label>
            <input id="dateTaken" v-model="metadata.dateTaken" type="date">
          </div>
          <div v-else>
            <div v-if="dateTaken" class="flex items-center space-x-3">
              <icon-tabler-calendar class="text-2xl text-gray-600" />
              <span class="text-gray-600">Taken on {{ dateTaken }}</span>
            </div>
          </div>
        </div>

        <div v-if="dateUploaded" class="flex items-center space-x-3">
          <icon-tabler-upload class="text-2xl text-gray-600" />
          <span class="text-gray-600">Uploaded on {{ dateUploaded }}</span>
        </div>

        <div v-if="metadata.exposureMode && metadata.exposureMode !== 'unknown'" class="flex items-center space-x-3">
          <icon-tabler-settings class="text-2xl text-gray-600" />
          <span class="text-gray-600">Mode: {{ metadata.exposureMode }}</span>
        </div>

        <div v-if="fStop" class="flex items-center space-x-3">
          <icon-tabler-aperture class="text-2xl text-gray-600" />
          <span class="text-gray-600">{{ fStop }}</span>
        </div>

        <div v-if="focalLength" class="flex items-center space-x-3">
          <icon-tabler-eye-pin class="text-2xl text-gray-600" />
          <span class="text-gray-600">{{ focalLength }}</span>
        </div>

        <div v-if="metadata.exposureTime && metadata.exposureTime !== 'unknown'" class="flex items-center space-x-3">
          <icon-tabler-clock class="text-2xl text-gray-600" />
          <span class="text-gray-600">{{ metadata.exposureTime }}</span>
        </div>

        <div v-if="metadata.iso && metadata.iso !== 'unknown'" class="flex items-center space-x-3">
          <icon-carbon-iso-outline class="text-2xl text-gray-600" />
          <span class="text-gray-600">{{ metadata.iso }}</span>
        </div>

        <div v-if="metadata.flashStatus && metadata.flashStatus !== 'unknown'" class="flex items-center space-x-3">
          <icon-tabler-bolt class="text-2xl text-gray-600" />
          <span class="text-gray-600">{{ metadata.flashStatus }}</span>
        </div>

        <div v-if="whiteBalance" class="flex items-center space-x-3">
          <icon-tabler-sun class="text-2xl text-gray-600" />
          <span class="text-gray-600">{{ whiteBalance }}</span>
        </div>
        <br>
        <div class="flex cursor-pointer items-center space-x-3" @click="loadOriginal()">
          <icon-tabler-arrow-autofit-up class="text-2xl" :class="loadOriginalIconColour" />
          <span class="text-gray-600">{{ loadOriginalText }}</span>
        </div>
        <div class="flex cursor-pointer items-center space-x-3" @click="downloadOriginal()">
          <icon-tabler-download class="text-2xl text-gray-600" />
          <span class="text-gray-600">Download original</span>
        </div>
        <div class="w-full flex flex-col gap-2 border-1 border-gray-400 border-solid p-4">
          <div class="text-center text-lg">
            This photo is in {{ imageAlbums.length }} albums
          </div>
          <div class="grid grid-cols-2 gap-2 gap-4 lg:grid-cols-4 md:grid-cols-3">
            <div v-for="(album, index) in imageAlbums" :key="index" class="flex flex-col gap-2">
              <img :src="getAlbumCoverSlugThumbnailAddress(album)" class="size-20" />
              <div class="overflow-auto text-center">
                {{ album.Name }}
              </div>
            </div>
            <div v-if="inEditingMode">
              <icon-tabler-square-plus class="size-20 text-gray-400" @click="addToAlbum" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
