<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { type Album, removeImageFromAlbum } from '../composables/albums'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import AlbumCoverSmall from './AlbumCoverSmall.vue'

const props = defineProps({
  imageSlug: { type: String, required: true },
  inEditingMode: { type: Boolean, default: false },
})

const emits = defineEmits(['addToAlbum'])

const router = useRouter()
const imageAlbums = ref<Album[]>([])

async function getAlbumsList() {
  imageAlbums.value = []
  const response = await backendFetchRequest(`links/image/${props.imageSlug}`)
  const albumSlugs: string[] = await response.json() || []
  albumSlugs.forEach(async (albumSlug) => {
    const response = await backendFetchRequest(`albums/${albumSlug}`)
    const album: Album = await response.json()
    imageAlbums.value.push(album)
  })
}

function navigateToAlbum(albumSlug: string) {
  router.push(`/albums/${albumSlug}`)
}

async function removeAlbumLink(album: Album) {
  console.log(`album: ${album.Slug}`)
  console.log(`image: ${props.imageSlug}`)
  await removeImageFromAlbum(album.Slug, props.imageSlug)
  getAlbumsList()
}

watch(
  () => props.imageSlug,
  () => {
    getAlbumsList()
  },
)

defineExpose({ getAlbumsList })

onMounted(() => {
  getAlbumsList()
})
</script>

<template>
  <div v-if="imageAlbums.length > 0 || inEditingMode" class="mt-4 max-w-120 flex flex-col gap-4 border-1 border-gray-400 rounded-sm border-solid p-4 dark:border-gray-600">
    <div class="text-left text-lg">
      This photo is in {{ imageAlbums.length }} {{ imageAlbums.length === 1 ? 'album' : 'albums' }}
    </div>
    <div v-if="inEditingMode || imageAlbums.length" class="grid grid-cols-2 gap-4 lg:grid-cols-4 md:grid-cols-3">
      <div v-for="(album, index) in imageAlbums" :key="index">
        <AlbumCoverSmall :album="album" :in-editing-mode="inEditingMode" :allow-delete="true" @image-click="navigateToAlbum(album.Slug)" @delete-click="removeAlbumLink" />
      </div>
      <div v-if="inEditingMode">
        <div class="hover:cursor-pointer" @click="emits('addToAlbum')">
          <icon-tabler-square-plus class="text-2xl" />
        </div>
      </div>
    </div>
  </div>
</template>
