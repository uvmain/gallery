<script setup lang="ts">
import { useSessionStorage } from '@vueuse/core'
import { onBeforeMount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const fileSelector = ref<HTMLInputElement | null>(null)
const title = ref()
const feedback = ref()

const userLoginState = useSessionStorage('login-state', false)
const router = useRouter()

function updateTitle() {
  const fileInput = fileSelector.value
  if (fileInput && fileInput.files && fileInput.files.length > 0) {
    let filename = fileInput.files[0].name
    filename = filename.replaceAll('-', ' ')
    filename = filename.replaceAll('_', ' ')
    filename = filename.substring(0, filename.lastIndexOf('.')) || filename
    title.value = filename.toLocaleLowerCase()
  }
}

async function uploadFile() {
  const fileInput = fileSelector.value
  if (!fileInput || !fileInput.files || fileInput.files.length === 0) {
    feedback.value = 'Please select a file before uploading.'
    return
  }

  const file = fileInput.files[0]
  const formData = new FormData()
  formData.append('file', file)
  formData.append('title', title.value)

  try {
    const response = await backendFetchRequest('upload', {
      body: formData,
      method: 'POST',
    })

    if (response.ok) {
      const slug = await response.json()
      router.push(`/${slug}`)
    }
    else {
      const errorText = await response.body
      feedback.value = `Upload failed: ${errorText}`
    }
  }
  catch (error) {
    feedback.value = `Error during upload: ${error}`
  }
}

onBeforeMount(() => {
  if (!userLoginState.value) {
    router.push('/')
  }
})
</script>

<template>
  <div>
    <Header />
    <div class="my-50 flex justify-center">
      <div class="flex flex-col gap-8">
        <div class="flex flex-col gap-4">
          <label for="fileSelector">Choose an image file:</label>
          <input
            id="fileSelector"
            ref="fileSelector"
            type="file"
            name="avatar"
            accept="image/avif, image/bmp, image/gif, image/jpg, image/jpeg, image/png, image/webp"
            class="border-none input"
            @change="updateTitle"
          />
        </div>
        <div class="flex flex-col gap-4">
          <label for="title">Title:</label>
          <input
            id="title"
            v-model="title"
            type="text"
            name="title"
            :multiple="false"
            class="border-solid input"
          />
        </div>
        <div>
          <button
            class="button"
            @click="uploadFile"
          >
            Upload
          </button>
        </div>
        <div v-if="feedback">
          {{ feedback }}
        </div>
      </div>
    </div>
  </div>
</template>
