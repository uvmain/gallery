<script setup lang="ts">
const postFailure = ref<string | null>()

async function post() {
  const imageMetadata: ImageMetadata = {
    fileName: "testing",
    title: "also testing",
    dateTaken: new Date(),
    dateUploaded: new Date(),
    cameraModel: "Olympus",
    lensModel: "60mm",
    aperture: "f/5.6",
    shutterSpeed: "1/2000",
    flashStatus: "Active, didn't flash",
    focusLength: "60mm",
    iso: "200",
    exposureMode: "Aperture Priority",
    whiteBalance: "Manual"
  }

  try {
    await $fetch("/api/image", {
      method: "POST",
      body: imageMetadata
    })
  }
  catch(error) {
    postFailure.value = `${error}`
    console.warn(error)
  }
}

onMounted(() => {
  post()
})
</script>

<template>
  <p v-if="postFailure">
    Failed to post: {{ postFailure }}
  </p>
  <p v-else>
    success
  </p>
</template>