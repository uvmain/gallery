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
  return `http://localhost:8080/api/optimised/${slug.value}`
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
    const response = await fetch(`http://localhost:8080/api/metadata/${slug.value}`)
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
  <div class="flex flex-col p-6 lg:flex-row">
    <img :src="optimisedPath" class="max-h-90vh max-w-70vw" />
    <div v-if="metadata" class="flex flex-col gap-y-2 p-6">
      <h1>
        {{ metadata.title }}
      </h1>
      <div v-if="camera" class="flex items-center space-x-2">
        <icon-tabler-camera class="size-4em text-black" />
        <div class="flex flex-col gap-2 text-lg">
          <span>{{ camera }}</span>
          <span>{{ lens }}</span>
        </div>
      </div>
      <div v-if="dateTaken" class="flex items-center space-x-2">
        <icon-tabler-camera-bolt class="text-2xl text-black" />
        <span>Taken on {{ dateTaken }}</span>
      </div>
      <div v-if="dateUploaded" class="flex items-center space-x-2">
        <icon-tabler-camera-up class="text-2xl text-black" />
        <span>Uploaded on {{ dateUploaded }}</span>
      </div>
      <div v-if="metadata.exposureMode" class="flex items-center space-x-2">
        <icon-tabler-camera-cog class="text-2xl text-black" />
        <span>Mode: {{ metadata.exposureMode }}</span>
      </div>
      <div v-if="fStop" class="flex items-center space-x-2">
        <icon-tabler-aperture class="text-2xl text-black" />
        <span>{{ fStop }}</span>
      </div>
      <div v-if="focalLength" class="flex items-center space-x-2">
        <icon-tabler-eye-pin class="text-2xl text-black" />
        <span>{{ focalLength }}</span>
      </div>
      <div v-if="metadata.exposureTime" class="flex items-center space-x-2">
        <icon-tabler-stopwatch class="text-2xl text-black" />
        <span>{{ metadata.exposureTime }}</span>
      </div>
      <div v-if="metadata.iso" class="flex items-center space-x-2">
        <icon-carbon-iso-outline class="text-2xl text-black" />
        <span>{{ metadata.iso }}</span>
      </div>
      <div v-if="metadata.flashStatus" class="flex items-center space-x-2">
        <icon-tabler-bolt class="text-2xl text-black" />
        <span>{{ metadata.flashStatus }}</span>
      </div>
      <div v-if="whiteBalance" class="flex items-center space-x-2">
        <icon-tabler-temperature-sun class="text-2xl text-black" />
        <span>{{ whiteBalance }}</span>
      </div>
    </div>
  </div>
</template>
