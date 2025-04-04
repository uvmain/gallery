<script setup lang="ts">
import type PhotoAlbums from '../components/PhotoAlbums.vue'
import type Tags from '../components/Tags.vue'
import type { ImageMetadata } from '../types/main'
import { onKeyStroke, useSessionStorage } from '@vueuse/core'
import dayjs from 'dayjs'
import { addImageToAlbum } from '../composables/albums'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import { getNextSlug, getPreviousSlug, getRandomSlug, getSlugPosition } from '../composables/slugs'

const route = useRoute()
const router = useRouter()

const slug = ref(route.params.imageSlug as string)
const slugPosition = ref<number>()
const prevSlug = ref<string>()
const nextSlug = ref<string>()
const imageSize = ref<string>('optimised')
const loadOriginalText = ref()
const metadata = ref<ImageMetadata | undefined>()
const tags = ref<typeof Tags>()
const inEditingMode = ref(false)
const selectedAlbums = ref<string[]>([])
const addToAlbumDialog = ref()
const photoalbums = ref<typeof PhotoAlbums>()

const userLoginState = useSessionStorage('login-state', false)

const imageSource = computed(() => {
  return `/api/${imageSize.value}/${slug.value}`
})

const fStop = computed(() => {
  const [first, second] = metadata.value ? metadata.value.fStop.split('/').map(Number) : [1, 1]
  const result = (first / second)
  return Number.isNaN(result) ? undefined : `f${result.toFixed(1)}`
})

const focalLength = computed(() => {
  const [first, second] = metadata.value ? metadata.value.focalLength.split('/').map(Number) : [1, 1]
  const result = (first / second)
  return Number.isNaN(result) ? undefined : `${result.toFixed(1)} mm`
})

const dateTaken = computed(() => {
  return metadata.value && metadata.value.dateTaken !== 'unknown' ? `Taken on ${dayjs(metadata.value.dateTaken).format('DD MMMM YYYY HH:mm:ss')}` : undefined
})

const dateUploaded = computed(() => {
  return metadata.value && metadata.value.dateUploaded !== 'unknown' ? `Uploaded on ${dayjs(metadata.value.dateUploaded).format('DD MMMM YYYY HH:mm:ss')}` : undefined
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
    result = (`${setting}` === `${temp}`) ? `${setting}` : `${setting}: ${temp}`
  }
  return result
})

const loadOriginalIconColour = computed(() => {
  return imageSize.value === 'optimised' ? '' : 'text-gray-400'
})

async function getMetadata() {
  try {
    slugPosition.value = await getSlugPosition(slug.value)
    prevSlug.value = await getPreviousSlug(slug.value)
    nextSlug.value = await getNextSlug(slug.value)
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

const patchMetadata = computed(() => {
  return {
    filePath: metadata.value?.filePath,
    fileName: metadata.value?.fileName,
    title: metadata.value?.title,
    dateTaken: metadata.value?.dateTaken,
    dateUploaded: metadata.value?.dateUploaded,
    cameraMake: metadata.value?.cameraMake,
    cameraModel: metadata.value?.cameraModel,
    lensMake: metadata.value?.lensMake,
    lensModel: metadata.value?.lensModel,
    fStop: metadata.value?.fStop,
    exposureTime: metadata.value?.exposureTime,
    flashStatus: metadata.value?.flashStatus,
    focalLength: metadata.value?.focalLength,
    iso: metadata.value?.iso,
    exposureMode: metadata.value?.exposureMode,
    whiteBalance: metadata.value?.whiteBalance,
    whiteBalanceMode: metadata.value?.whiteBalance,
  }
})

async function saveMetadata() {
  if (userLoginState.value) {
    inEditingMode.value = false
    try {
      const options = {
        body: JSON.stringify(patchMetadata.value),
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
      }
      await backendFetchRequest(`metadata/${slug.value}`, options)
      await getMetadata()
      await tags.value?.getTags()
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
  addToAlbumDialog.value?.show()
}

function hideAddToAlbumDialog() {
  addToAlbumDialog.value?.hide()
}

async function ConfirmAddToAlbum(albumSlugArray: string[]) {
  addToAlbumDialog.value?.hide()
  for (const albumSlug of albumSlugArray) {
    await addImageToAlbum(albumSlug, slug.value)
  }
  photoalbums.value?.getAlbumsList()
}

async function deleteImage() {
  if (userLoginState.value) {
    try {
      const options = {
        body: JSON.stringify({ slug: slug.value }),
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
      }
      const response = await backendFetchRequest(`slugs/${slug.value}`, options)
      if (response.status === 200) {
        router.push('/')
      }
    }
    catch (error) {
      console.error('Failed to delete image:', error)
    }
  }
  else {
    console.warn('Not authorised to delete image, please log in')
  }
}

watch(
  () => route.params.imageSlug,
  () => {
    slug.value = route.params.imageSlug as string
    getMetadata()
  },
)

onKeyStroke('ArrowLeft', (e) => {
  if (inEditingMode.value) {
    return
  }
  e.preventDefault()
  router.push(`/${prevSlug.value}`)
})

onKeyStroke('ArrowRight', (e) => {
  if (inEditingMode.value) {
    return
  }
  e.preventDefault()
  router.push(`/${nextSlug.value}`)
})

onKeyStroke('ArrowDown', async (e) => {
  if (inEditingMode.value) {
    return
  }
  e.preventDefault()
  const slug = await getRandomSlug()
  router.push(`/${slug}`)
})

onKeyStroke('ArrowUp', async (e) => {
  if (inEditingMode.value) {
    return
  }
  e.preventDefault()
  router.go(-1)
})

onBeforeMount(async () => {
  await getMetadata()
})
</script>

<template>
  <div class="min-h-screen">
    <Header :show-edit="!inEditingMode" @edit="enableEditing">
      <template #2>
        <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="deleteImage">
          <icon-tabler-trash class="text-2xl text-red-700" />
        </div>
        <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="disableEditing">
          <icon-tabler-edit-off class="text-2xl text-orange-700" />
        </div>
        <div v-if="inEditingMode" class="p-2 hover:cursor-pointer" @click="saveMetadata">
          <icon-tabler-checkbox class="text-2xl text-green-700" />
        </div>
      </template>
    </Header>
    <div class="flex flex-col justify-center gap-x-8 gap-y-4 p-2 lg:flex-row">
      <div class="mx-auto lg:mx-0 lg:pt-4">
        <ZoomableImage :image-source="imageSource" :width="metadata?.width" :height="metadata?.height" />
      </div>
      <div v-if="metadata" class="flex flex-col gap-3 text-sm lg:max-w-1/3">
        <TextInput v-model="metadata.title" :in-editing-mode="inEditingMode" class="py-2 text-xl font-semibold lg:text-3xl" @keypress.enter="saveMetadata" />

        <div v-if="camera" class="flex items-center space-x-3">
          <div class="group">
            <icon-tabler-camera class="text-2xl" />
            <span class="tooltip">
              Camera Model
            </span>
          </div>
          <div class="flex flex-col gap-1 text-base">
            <span>{{ camera }}</span>
            <span>{{ lens }}</span>
          </div>
        </div>

        <div>
          <div class="flex gap-x-2">
            <TooltipIcon tooltip-text="Date Taken">
              <icon-tabler-calendar class="text-xl" />
            </TooltipIcon>
            <DateInput v-model="metadata.dateTaken" :in-editing-mode="inEditingMode" :display-value="dateTaken" class="py-2" @keypress.enter="saveMetadata" />
          </div>
        </div>

        <TooltipIcon v-if="dateUploaded" tooltip-text="Date Uploaded" :content="dateUploaded">
          <icon-tabler-upload class="text-xl" />
        </TooltipIcon>

        <TooltipIcon v-if="metadata.exposureMode && metadata.exposureMode !== 'unknown'" tooltip-text="Exposure Mode" :content="metadata.exposureMode">
          <icon-tabler-settings class="text-xl" />
        </TooltipIcon>

        <TooltipIcon v-if="fStop" tooltip-text="fStop" :content="fStop">
          <icon-tabler-aperture class="text-xl" />
        </TooltipIcon>

        <TooltipIcon v-if="focalLength" tooltip-text="Focal Length" :content="focalLength">
          <icon-tabler-eye-pin class="text-xl" />
        </TooltipIcon>

        <TooltipIcon v-if="metadata.exposureTime && metadata.exposureTime !== 'unknown'" tooltip-text="Shutter Speed" :content="metadata.exposureTime">
          <icon-tabler-clock class="text-xl" />
        </TooltipIcon>

        <TooltipIcon v-if="metadata.iso && metadata.iso !== 'unknown'" tooltip-text="ISO" :content="metadata.iso">
          <icon-carbon-iso-outline class="text-xl" />
        </TooltipIcon>

        <TooltipIcon v-if="metadata.flashStatus && metadata.flashStatus !== 'unknown'" tooltip-text="Flash Status" :content="metadata.flashStatus">
          <icon-tabler-bolt class="text-xl" />
        </TooltipIcon>

        <TooltipIcon v-if="whiteBalance" tooltip-text="White Balance" :content="whiteBalance">
          <icon-tabler-sun class="text-xl" />
        </TooltipIcon>

        <br>
        <div class="flex cursor-pointer items-center space-x-3" @click="loadOriginal()">
          <icon-tabler-arrow-autofit-up id="load-original" class="text-xl" :class="loadOriginalIconColour" />
          <label for="load-original">{{ loadOriginalText }}</label>
        </div>
        <div class="flex cursor-pointer items-center space-x-3" @click="downloadOriginal()">
          <icon-tabler-download id="download-original" class="text-xl" />
          <label for="download-original">Download original</label>
        </div>
        <PhotoAlbums ref="photoalbums" v-model:in-editing-mode="inEditingMode" :image-slug="slug" @add-to-album="addToAlbum()" />
        <Tags ref="tags" v-model:in-editing-mode="inEditingMode" :image-slug="slug" />
      </div>
    </div>
    <Dialog v-if="metadata" ref="addToAlbumDialog" :close-button="false" class="size-90%" @keydown.escape="hideAddToAlbumDialog()">
      <div class="flex flex-col items-center justify-center gap-4 lg:flex-row">
        <img
          :src="imageSource"
          :alt="slug"
          loading="lazy"
          onerror="this.onerror=null;this.src='/default-image.jpg';"
          class="h-40 w-80 cursor-pointer border-2 border-white border-solid object-cover dark:border-neutral-500"
        />
        <div class="text-lg">
          {{ metadata.title }}
        </div>
      </div>
      <AlbumSelector v-model:selected-albums="selectedAlbums" :single-select="false" :in-editing-mode="inEditingMode" @cancel="hideAddToAlbumDialog()" @confirm="ConfirmAddToAlbum" />
    </Dialog>
  </div>
</template>
