<script setup lang="ts">
import { ref } from 'vue'

const { loggedIn } = useUserSession()
const slugs = ref<string[]>([])
const showModal = ref(false)
const selectedImage = ref<string | null>(null)

async function getSlugs() {
  if (loggedIn.value) {
    try {
      const response = await $fetch<string[]>("/api/slugs")
        .catch((error) => {
          console.error(`Failed to fetch data: ${JSON.stringify(error.data)}`)
        })
      if (response) {
        slugs.value = response
      }
    }
    catch (error) {
      console.error('Failed to fetch thumbnails:', error)
    }
  }
}

function openModal(thumbnailPath: string) {
  selectedImage.value = thumbnailPath
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  selectedImage.value = null
}

function getThumbnailPath(slug: string) {
  return `/api/thumbnail/${slug}`
}

function getImagePath(slug: string) {
  return `/api/images/${slug}`
}

onBeforeMount(async () => {
  getSlugs()
})
</script>

<template>
  <div class="flex flex-col gap-2 items-center justify-center">
    <div v-if="loggedIn" class="flex flex-wrap gap-2 justify-start lg:max-w-8/10">
      <div v-for="(slug, index) in slugs" :key="index">
        <img :src="getThumbnailPath(slug)" :alt="slug" class="h-20vh object-cover cursor-pointer" @click="openModal(getImagePath(slug))" >
      </div>
        <!-- Modal -->
      <div v-if="showModal" class="inset-0 fixed bg-black bg-opacity-75 flex items-center justify-center z-50">
        <div class="relative bg-white p-4 rounded-lg">
          <button class="absolute top-2 right-2 text-white bg-red-500 rounded-lg size-6" @click="closeModal">X</button>
          <img v-if="selectedImage" :src="selectedImage" alt="Selected Image" class="max-w-full max-h-screen" >
        </div>
      </div>
    </div>
    <p v-else>
      Please log in
    </p>
  </div>
</template>