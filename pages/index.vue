<script setup lang="ts">
const { loggedIn } = useUserSession()

const imageData = ref()
const imagesLastUpdated = useState('imagesLastUpdated')

// async function getFullSizeImage() {
//   if (loggedIn.value) {
//     const response = await fetch("/api/thumbnail/20240906-house-moth.jpg")
    
//     if (response.ok) {
//       const blob = await response.blob()
//       imageData.value = URL.createObjectURL(blob)
//     }
//     else {
//       console.error("Failed to fetch image")
//     }
//   }
// }

// watch(loggedIn, (newValue) => {
//   console.info(`loggedIn: ${newValue}`)
//   if (newValue == true) {
//     getFullSizeImage()
//   }
// })

onBeforeMount(async () => {
  // getFullSizeImage()
})
</script>

<template>
  <div class="flex flex-col gap-2 items-center justify-center">
    <div v-if="loggedIn">
      {{ imagesLastUpdated }}
      <img v-if="imageData" :src="imageData" >
    </div>
    <p v-else>
      Please log in
    </p>
  </div>
</template>