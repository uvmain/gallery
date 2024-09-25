<script setup lang="ts">
const postStatus = ref()

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
    const status = await $fetch("/api/image", {
      method: "POST",
      body: imageMetadata
    })
    postStatus.value = status
  }
  catch(error) {
    postStatus.value = `${error}`
    console.warn(error)
  }
}

onMounted(() => {
  post()
})
</script>

<template>
  <div>
    {{ postStatus }}
  </div>
</template>