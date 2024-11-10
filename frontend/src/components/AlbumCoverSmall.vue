<script setup lang="ts">
import type { Album } from '../composables/albums'
import { onBeforeMount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getAlbumCoverSlugThumbnailAddress } from '../composables/albums'

const props = withDefaults(defineProps<{
  album: Album
  showName?: boolean
}>(), {
  showName: true,
})

const albumThumbnailAddress = ref()

const router = useRouter()

function navigateToAlbum() {
  router.push(`/albums/${props.album.Slug}`)
}

onBeforeMount(async () => {
  albumThumbnailAddress.value = await getAlbumCoverSlugThumbnailAddress(props.album)
})
</script>

<template>
  <div>
    <img :src="albumThumbnailAddress" onerror="this.onerror=null;this.src='/default-image.jpg';" class="size-20 border-2 border-white border-solid hover:cursor-pointer" @click="navigateToAlbum" />
    <div v-if="showName" class="max-w-20 overflow-hidden text-center text-gray-600">
      {{ album.Name }}
    </div>
  </div>
</template>
