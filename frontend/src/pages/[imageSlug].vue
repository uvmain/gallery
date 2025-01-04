<script setup lang="ts">
import type { ImageMetadata } from '../composables/imageMetadataInterface'
import { useSessionStorage } from '@vueuse/core'
import dayjs from 'dayjs'
import { computed, onBeforeMount, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const route = useRoute()
const router = useRouter()

const slug = ref(route.params.imageSlug as string)
const imageSize = ref<string>('optimised')
const loadOriginalText = ref()
const metadata = ref<ImageMetadata | undefined>()
const inEditingMode = ref(false)

const userLoginState = useSessionStorage('login-state', false)
const selectedImage = useSessionStorage('selected-image', '')

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
  return imageSize.value === 'optimised' ? '' : 'text-gray-400'
})

async function getMetadata() {
  try {
    const response = await backendFetchRequest(`metadata/${slug.value}`)
    metadata.value = await response.json() as ImageMetadata
    imageSize.value = 'optimised'
    loadOriginalText.value = 'Load Original'
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
  selectedImage.value = route.params.imageSlug as string
  router.push('/albums')
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
  <div class="min-h-screen">
    <Header :show-edit="!inEditingMode" @edit="enableEditing">
      <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="disableEditing">
        <icon-tabler-edit-off class="text-2xl text-red-700" />
      </div>
      <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="saveMetadata">
        <icon-tabler-checkbox class="text-2xl text-green-700" />
      </div>
    </Header>
    <div id="main" class="flex flex-row justify-center gap-6 p-6">
      <img v-if="imageSource" :src="imageSource" class="max-h-80vh max-w-70vw border-6 border-white border-solid dark:border-neutral-500" />
      <div v-if="metadata" class="flex flex-col gap-4 p-6 text-sm lg:max-w-1/3">
        <div>
          <div v-if="inEditingMode">
            <label for="imageTitle" class="mb-4 text-2xl">
              Title:
            </label>
            <input id="imageTitle" v-model="metadata.title" type="text" @keypress.enter="saveMetadata">
          </div>
          <h1 v-else class="mb-4 text-2xl font-semibold">
            {{ metadata.title }}
          </h1>
        </div>

        <div v-if="camera" class="flex items-center space-x-3">
          <icon-tabler-camera class="text-3xl" />
          <div class="flex flex-col gap-1 text-base">
            <span>{{ camera }}</span>
            <span>{{ lens }}</span>
          </div>
        </div>

        <div>
          <div v-if="inEditingMode">
            <label for="dateTaken" class="mb-4 text-2xl">
              Date taken:
            </label>
            <input id="dateTaken" v-model="metadata.dateTaken" type="datetime-local">
          </div>
          <div v-else>
            <div v-if="dateTaken" class="flex items-center space-x-3">
              <icon-tabler-calendar class="text-2xl" />
              <span>Taken on {{ dateTaken }}</span>
            </div>
          </div>
        </div>

        <div v-if="dateUploaded" class="flex items-center space-x-3">
          <icon-tabler-upload class="text-2xl" />
          <span>Uploaded on {{ dateUploaded }}</span>
        </div>

        <div v-if="metadata.exposureMode && metadata.exposureMode !== 'unknown'" class="flex items-center space-x-3">
          <icon-tabler-settings class="text-2xl" />
          <span>Mode: {{ metadata.exposureMode }}</span>
        </div>

        <div v-if="fStop" class="flex items-center space-x-3">
          <icon-tabler-aperture class="text-2xl" />
          <span>{{ fStop }}</span>
        </div>

        <div v-if="focalLength" class="flex items-center space-x-3">
          <icon-tabler-eye-pin class="text-2xl" />
          <span>{{ focalLength }}</span>
        </div>

        <div v-if="metadata.exposureTime && metadata.exposureTime !== 'unknown'" class="flex items-center space-x-3">
          <icon-tabler-clock class="text-2xl" />
          <span>{{ metadata.exposureTime }}</span>
        </div>

        <div v-if="metadata.iso && metadata.iso !== 'unknown'" class="flex items-center space-x-3">
          <icon-carbon-iso-outline class="text-2xl" />
          <span>{{ metadata.iso }}</span>
        </div>

        <div v-if="metadata.flashStatus && metadata.flashStatus !== 'unknown'" class="flex items-center space-x-3">
          <icon-tabler-bolt class="text-2xl" />
          <span>{{ metadata.flashStatus }}</span>
        </div>

        <div v-if="whiteBalance" class="flex items-center space-x-3">
          <icon-tabler-sun class="text-2xl" />
          <span>{{ whiteBalance }}</span>
        </div>
        <br>
        <div class="flex cursor-pointer items-center space-x-3" @click="loadOriginal()">
          <icon-tabler-arrow-autofit-up id="load-original" class="text-2xl" :class="loadOriginalIconColour" />
          <label for="load-original">{{ loadOriginalText }}</label>
        </div>
        <div class="flex cursor-pointer items-center space-x-3" @click="downloadOriginal()">
          <icon-tabler-download id="download-original" class="text-2xl" />
          <label for="download-original">Download original</label>
        </div>
        <PhotoAlbums v-model:in-editing-mode="inEditingMode" :image-slug="slug" @add-to-album="addToAlbum()" />
      </div>
    </div>
  </div>
</template>
