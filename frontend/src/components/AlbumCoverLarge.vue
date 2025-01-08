<script setup lang="ts">
import type { Album } from '../composables/albums'
import { onBeforeMount, ref } from 'vue'
import { getAlbumCoverSlugThumbnailAddress } from '../composables/albums'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import { niceDate } from '../composables/logic'

const props = defineProps<{
  album: Album
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
      <hr class="mx-auto my-2px h-px max-w-60% border-0 bg-gray-400 opacity-60">
      <hr class="mx-auto my-2px h-px max-w-70% border-0 bg-gray-400 opacity-80">
      <hr class="mx-auto my-2px h-px max-w-80% border-0 bg-gray-400">
      <img :src="albumThumbnailAddress" onerror="this.onerror=null;this.src='/default-image.jpg';" class="size-60 border-4 border-white border-solid object-cover dark:border-neutral-500" @click="emits('navigate', album.Slug)" />
      <div class="absolute right-4 top-3 p-2 hover:cursor-pointer" @click="emits('trash', album)">
        <icon-tabler-trash-x class="text-xl hover:size-115%" />
      </div>
      <div class="grad absolute bottom-2 left-2 mb-1 w-auto flex flex-col gap-2 rounded p-2 text-white">
        <div class="[text-shadow:_0_0px_4px_rgb(0_0_0_/_0.8)] text-lg font-semibold">
          {{ album.Name }}
        </div>
        <div class="[text-shadow:_0_0px_4px_rgb(0_0_0_/_0.8)]">
          {{ imageCount }} photos
        </div>
        <div class="[text-shadow:_0_0px_4px_rgb(0_0_0_/_0.8)]">
          Created: {{ niceDate(album.DateCreated) }}
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.grad {
  background: radial-gradient(circle at bottom left, rgb(0, 0, 0, 0.5) 0, rgba(255, 255, 255, 0) 60%);
}
</style>
