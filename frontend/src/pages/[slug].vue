<script setup lang="ts">
import dayjs from 'dayjs'
import { computed, onBeforeMount, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const slug = ref(route.params.slug as string)

interface ImageMetadata {
  filePath: string
  fileName: string
  title: string
  dateTaken: string
  dateUploaded: string
  cameraMake: string
  cameraModel: string
  lensMake: string
  lensModel: string
  fStop: string
  exposureTime: string
  flashStatus: string
  focalLength: string
  iso: string
  exposureMode: string
  whiteBalance: string
  whiteBalanceMode: string
  albums: '[]'
}

const metadata = ref<ImageMetadata | undefined>()

const optimisedPath = computed(() => {
  return `/api/optimised/${slug.value}`
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

async function getMetadata() {
  try {
    const response = await fetch(`/api/metadata/${slug.value}`)
    metadata.value = await response.json() as ImageMetadata
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
  }
}

onBeforeMount(() => {
  getMetadata()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <div id="main" class="flex flex-row justify-center gap-8 p-6">
      <!-- Image Section -->
      <div class="border-6 border-white border-solid">
        <img :src="optimisedPath" class="max-h-90vh max-w-70vw" />
      </div>

      <!-- EXIF Data Section -->
      <div v-if="metadata" class="flex flex-col gap-6 p-6 text-sm lg:max-w-1/3">
        <h1 class="mb-4 text-2xl text-gray-700 font-semibold">
          {{ metadata.title }}
        </h1>

        <div v-if="camera" class="flex items-center space-x-3">
          <icon-tabler-camera class="text-3xl text-gray-600" />
          <div class="flex flex-col gap-1 text-base">
            <span class="text-gray-700 font-medium">{{ camera }}</span>
            <span class="text-gray-600">{{ lens }}</span>
          </div>
        </div>

        <div v-if="dateTaken" class="flex items-center space-x-3">
          <icon-tabler-calendar class="text-2xl text-gray-600" />
          <span class="text-gray-600">Taken on {{ dateTaken }}</span>
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
      </div>
    </div>
  </div>
</template>
