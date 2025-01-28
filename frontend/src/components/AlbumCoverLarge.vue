<script setup lang="ts">
import type { Album } from '../composables/albums'
import { getAlbumCoverSlugThumbnailAddress } from '../composables/albums'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const props = defineProps<{
  album: Album
  inEditMode: boolean
}>()

const emits = defineEmits(['trash', 'navigate'])

const albumThumbnailAddress = ref()
const imageCount = ref(0)

async function getImageCount() {
  const response = await backendFetchRequest(`links/album/${props.album.Slug}`)
  const imageSlugs: string[] = await response.json() || []
  imageCount.value = imageSlugs.length
}

onBeforeMount(async () => {
  albumThumbnailAddress.value = await getAlbumCoverSlugThumbnailAddress(props.album)
  getImageCount()
})
</script>

<template>
  <div>
    <div>
      <hr class="mx-auto my-2px hidden h-px max-w-60% border-0 bg-gray-400 opacity-60 lg:block">
      <hr class="mx-auto my-2px hidden h-px max-w-70% border-0 bg-gray-400 opacity-80 lg:block">
      <hr class="mx-auto my-2px hidden h-px max-w-80% border-0 bg-gray-400 lg:block">
      <img
        :src="albumThumbnailAddress"
        onerror="this.onerror=null;this.src='/default-image.jpg';"
        class="w-full border-4 border-white border-solid object-cover lg:size-60 dark:border-neutral-500"
        loading="lazy"
        @click="emits('navigate', album.Slug)"
      />
      <div v-if="inEditMode" class="grad absolute right-4 top-3 p-2 hover:cursor-pointer" @click="emits('trash', album)">
        <icon-tabler-trash-x class="text-xl text-white hover:text-green" />
      </div>
      <div class="absolute bottom-1 left-1 mb-1 w-8/10 flex flex-col gap-2 rounded-sm from-black from-opacity-50 to-opacity-0 bg-gradient-to-r p-2 text-white">
        <div class="[text-shadow:_0_0px_4px_rgb(0_0_0_/_0.8)] text-lg font-semibold">
          {{ album.Name }}
        </div>
        <div class="[text-shadow:_0_0px_4px_rgb(0_0_0_/_0.8)]">
          {{ imageCount }} photos
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.grad {
  background: radial-gradient(circle at top right, rgb(0, 0, 0, 0.5) 0, rgba(255, 255, 255, 0) 60%);
}
</style>
