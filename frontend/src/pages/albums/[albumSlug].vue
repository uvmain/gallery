<script setup lang="ts">
import { computed, onBeforeMount, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { backendFetchRequest, getServerUrl } from '../../composables/fetchFromBackend'
import { useStorage } from '@vueuse/core'

const route = useRoute()
const router = useRouter()

const serverBaseUrl = ref()
const albumData = ref()
const albumSlug = ref(route.params.albumSlug as string)
const userLoginState = useStorage('login-state', false)

const iconColour = computed(() => {
  return userLoginState.value ? 'text-gray-600' : 'text-gray-400'
})

async function getAlbumData() {
  const response = await backendFetchRequest(`albums/${albumSlug.value}`)
  albumData.value = await response.json()
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

const imageSource = computed(() => {
  const imageSlug = albumData.value.CoverSlug
  return `${serverBaseUrl.value}/api/optimised/${imageSlug}`
})

onBeforeMount(async () => {
  serverBaseUrl.value = await getServerUrl()
  getAlbumData()
})
</script>

<template>
  <div class="min-h-screen bg-gray-300">
    <Header bg="300" :show-edit="true">
      <div class="p-2 hover:cursor-pointer" @click="deleteAlbum">
        <icon-tabler-trash-x class="text-2xl" :class="iconColour" />
      </div>
    </Header>
    <div v-if="albumData" class="flex flex-row justify-center gap-8 p-6">
      {{ albumData }}
    </div>
    <div class="relative mx-auto flex flex-wrap gap-x-1 lg:max-w-8/10">
      <img :src="imageSource" :alt="albumSlug" class="h-50 w-full object-cover">
      <div class="absolute top-1/5 h-30 w-full bg-white opacity-90 blur-2xl">
      </div>
      <div class="absolute left-1/2 top-1/3 text-center text-3xl text-gray-600">
        {{ albumData.Name }}
      </div>
    </div>
  </div>
</template>

<style lang="css" scoped>
.bgImgCenter{
    background-image: url('imagePath');
    background-repeat: no-repeat;
    background-position: center;
    position: relative;
}
</style>
