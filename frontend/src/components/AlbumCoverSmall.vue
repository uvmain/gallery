<script setup lang="ts">
import type { Album } from '../composables/albums'
import { onBeforeMount, ref } from 'vue'
import { getAlbumCoverSlugThumbnailAddress } from '../composables/albums'

const props = defineProps<{
  album: Album
}>()

const albumThumbnailAddress = ref()

onBeforeMount(async () => {
  albumThumbnailAddress.value = await getAlbumCoverSlugThumbnailAddress(props.album)
})
</script>

<template>
  <div>
    <img :src="albumThumbnailAddress" class="size-20 border-2 border-white border-solid" />
    <div class="max-w-20 overflow-hidden text-center">
      {{ album.Name }}
    </div>
  </div>
</template>
