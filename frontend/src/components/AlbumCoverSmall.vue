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
  <div class="flex flex-wrap gap-2 justify-center relative">
    <div v-if="inEditingMode && allowDelete">
      <icon-tabler-circle-dashed-minus class="text-2xl text-red rounded-2xl bg-white bg-op-60 right-4 absolute z-100 hover-bg-op-100 hover:cursor-pointer -top-1" @click="emits('deleteClick', album)" />
    </div>
    <img
      :src="albumThumbnailAddress"
      loading="lazy"
      onerror="this.onerror=null;this.src='/default-image.jpg';"
      class="border-2 border-white border-solid size-20 dark:border-neutral-500 hover:cursor-pointer"
      @click="emits('imageClick', album)"
    />
    <div v-if="showName" class="text-center max-w-24">
      {{ album.Name }}
    </div>
  </div>
</template>
