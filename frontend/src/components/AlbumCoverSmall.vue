<script setup lang="ts">
import type { Album } from '../composables/albums'
import { getAlbumCoverSlugThumbnailAddress } from '../composables/albums'

const props = withDefaults(defineProps<{
  album: Album
  showName?: boolean
  inEditingMode?: boolean
  allowDelete?: boolean
}>(), {
  showName: true,
})

const emits = defineEmits(['deleteClick', 'imageClick'])

const albumThumbnailAddress = ref()

onBeforeMount(async () => {
  albumThumbnailAddress.value = await getAlbumCoverSlugThumbnailAddress(props.album)
})
</script>

<template>
  <div class="relative flex flex-wrap justify-center gap-2">
    <div v-if="inEditingMode && allowDelete">
      <icon-tabler-circle-dashed-minus class="absolute right-4 z-100 rounded-2xl bg-white bg-op-60 text-2xl text-red -top-1 hover:cursor-pointer hover-bg-op-100" @click="emits('deleteClick', album)" />
    </div>
    <img
      :src="albumThumbnailAddress"
      loading="lazy"
      onerror="this.onerror=null;this.src='/default-image.jpg';"
      class="size-20 border-2 border-white border-solid hover:cursor-pointer dark:border-neutral-500"
      @click="emits('imageClick', album)"
    />
    <div v-if="showName" class="max-w-24 text-center">
      {{ album.Name }}
    </div>
  </div>
</template>
