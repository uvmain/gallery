<script setup lang="ts">
const { loggedIn } = useUserSession()

const thumbnailPaths = ref<string[]>([])

async function getThumbnailPaths() {
  if (loggedIn.value) {
    try {
      const response = await $fetch<string[]>("/api/thumbnails")
        .catch((error) => {
          console.error(`Failed to fetch data: ${JSON.stringify(error.data)}`)
        })
      if (response) {
        thumbnailPaths.value = response
      }
    }
    catch (error) {
      console.error('Failed to fetch thumbnails:', error)
    }
  }
}

onBeforeMount(async () => {
  getThumbnailPaths()
})
</script>

<template>
  <div class="flex flex-col gap-2 items-center justify-center">
    <div v-if="loggedIn" class="flex flex-wrap gap-2 justify-start lg:max-w-8/10 p-4">
      <div v-for="(thumbnailPath, index) in thumbnailPaths" :key="index">
        <img :src="thumbnailPath" :alt="thumbnailPath" class="h-20vh object-cover">
      </div>
    </div>
    <p v-else>
      Please log in
    </p>
  </div>
</template>