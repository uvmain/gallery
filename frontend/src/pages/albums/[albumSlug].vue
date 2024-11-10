<script setup lang="ts">
import { useStorage } from '@vueuse/core'
import { computed, onBeforeMount, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { backendFetchRequest } from '../../composables/fetchFromBackend'
import { getThumbnailPath, niceDate } from '../../composables/logic'

const route = useRoute()
const router = useRouter()

const albumData = ref()
const albumLinks = ref<string[]>([])
const albumSlug = ref(route.params.albumSlug as string)
const userLoginState = useStorage('login-state', false)

const iconColour = computed(() => {
  return userLoginState.value ? 'text-gray-600' : 'text-gray-400'
})

async function getAlbumData() {
  const response = await backendFetchRequest(`albums/${albumSlug.value}`)
  albumData.value = await response.json()
}

async function getLinkData() {
  const response = await backendFetchRequest(`links/album/${albumSlug.value}`)
  albumLinks.value = await response.json()
}

async function deleteAlbum() {
  if (!userLoginState.value) {
    return
  }
  const options = {
    method: 'DELETE',
  }
  try {
    await backendFetchRequest(`albums/${albumSlug.value}`, options)
    router.push('/albums')
  }
  catch (error) {
    console.error('Failed to delete album:', error)
  }
}

onBeforeMount(async () => {
  getAlbumData()
  getLinkData()
})
</script>

<template>
  <div class="min-h-screen bg-gray-300">
    <Header bg="300" :show-edit="true">
      <div class="p-2 hover:cursor-pointer" @click="deleteAlbum">
        <icon-tabler-trash-x class="text-2xl" :class="iconColour" />
      </div>
    </Header>
    <div v-if="albumData" class="flex flex-row items-center justify-center gap-6 p-6 lg:max-w-8/10">
      <img
        :src="getThumbnailPath(albumData.CoverSlug)"
        :alt="albumData.CoverSlug"
        onerror="this.onerror=null;this.src='/default-image.jpg';"
        class="h-40 w-80 cursor-pointer border-2 border-white border-solid object-cover"
      />
      <div class="flex flex-col gap-2">
        <div class="text-2xl text-gray-600">
          {{ albumData.Name }}
        </div>
        <div class="text-gray-600">
          Created: {{ niceDate(albumData.DateCreated) }}
        </div>
        <div v-if="albumLinks" class="text-gray-600">
          {{ albumLinks.length }} photos
        </div>
      </div>
    </div>
    <div id="main" class="grid grid-cols-2 mx-auto gap-8 p-6 lg:grid-cols-7 md:grid-cols-4 lg:max-w-8/10">
      <div v-for="(imageSlug, index) in albumLinks" :key="index" class="relative">
        <img
          :src="getThumbnailPath(imageSlug)"
          :alt="imageSlug"
          class="h-full max-h-25vh max-w-40vw min-h-20vh w-full cursor-pointer object-cover"
          @click="router.push(`/${imageSlug}`)"
        />
      </div>
    </div>
  </div>
</template>
