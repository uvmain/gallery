<script setup lang="ts">
import type { Album } from '../../composables/albums'
import dayjs from 'dayjs'
import { onBeforeMount, ref } from 'vue'
import { getServerUrl } from '../../composables/getServerBaseUrl'

const serverBaseUrl = ref()
const albums = ref()

function addAlbum() {
  console.log('test')
}

async function getAlbums() {
  try {
    const response = await fetch(`${serverBaseUrl.value}/api/albums`)
    albums.value = await response.json() as Album
  }
  catch (error) {
    console.error('Failed to fetch Albums:', error)
  }
}

function getImageSource(slug = 'none') {
  return slug === 'none' ? '/default-image.jpg' : `${serverBaseUrl.value}/api/optimised/${slug}`
}

function niceDate(dateString: string) {
  return dayjs(dateString).format('DD/MM/YYYY')
}

function trashAlbum(albumName: string) {
  console.log(`trashing album ${albumName}`)
}

onBeforeMount(async () => {
  serverBaseUrl.value = await getServerUrl()
  await getAlbums()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <Header @add="addAlbum" />
    <div id="main" class="grid grid-cols-2 mx-auto gap-8 p-6 lg:grid-cols-7 md:grid-cols-4 lg:max-w-8/10">
      <div v-for="(album, index) in albums" :key="index" class="relative">
        <hr class="mx-auto my-2px h-px max-w-60% border-0 bg-gray-400 opacity-60">
        <hr class="mx-auto my-2px h-px max-w-70% border-0 bg-gray-400 opacity-80">
        <hr class="mx-auto my-2px h-px max-w-80% border-0 bg-gray-400">
        <img :src="getImageSource(album.CoverSlug)" onerror="this.onerror=null;this.src='/default-image.jpg';" class="size-60 border-4 border-white border-solid object-cover" />
        <div class="absolute right-4 top-2 p-2 hover:cursor-pointer" @click="trashAlbum(album.Name)">
          <icon-tabler-trash-x class="text-xl text-white hover:size-115%" />
        </div>
        <div class="absolute bottom-2 left-2 w-auto flex flex-col gap-2 p-2 text-white">
          <div class="[text-shadow:_0_0px_4px_rgb(0_0_0_/_0.8)] text-lg font-semibold">
            {{ album.Name }}
          </div>
          <div class="[text-shadow:_0_0px_4px_rgb(0_0_0_/_0.8)]">
            Created: {{ niceDate(album.DateCreated) }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
