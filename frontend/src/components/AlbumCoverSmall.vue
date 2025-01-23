<script setup lang="ts">
import type { Album } from '../composables/albums'
import { onBeforeMount, ref } from 'vue'
import { getAlbumCoverSlugThumbnailAddress } from '../composables/albums'

const props = withDefaults(defineProps<{
  album: Album
  showName?: boolean
  inEditingMode?: boolean
  allowDelete?: boolean
}>(), {
  showName: true,
})

const emits = defineEmits(['delete-click', 'image-click'])

const albumThumbnailAddress = ref()

onBeforeMount(async () => {
  albumThumbnailAddress.value = await getAlbumCoverSlugThumbnailAddress(props.album)
})
</script>

<template>
  <div class="relative">
    <div v-if="inEditingMode && allowDelete">
      <icon-tabler-circle-dashed-minus class="absolute right-4 z-100 rounded-2xl bg-white bg-op-60 text-2xl text-red -top-1 hover:cursor-pointer hover-bg-op-100" @click="emits('delete-click', album)" />
    </div>
    <img
      :src="albumThumbnailAddress"
      onerror="this.onerror=null;this.src='/default-image.jpg';"
      class="size-20 border-2 border-white border-solid hover:cursor-pointer dark:border-neutral-500"
      @click="emits('image-click', album)"
    />
    <div v-if="showName" class="max-w-20 overflow-hidden text-center">
      {{ album.Name }}
    </div>
  </div>
</template>
