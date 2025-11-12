<script setup lang="ts">
import { useSessionStorage } from '@vueuse/core'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const fileSelector = ref<HTMLInputElement | null>(null)
const title = ref()
const feedback = ref()
const image = ref()
const uploadDisabled = ref(false)

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
    uploadDisabled.value = true
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
  finally {
    uploadDisabled.value = false
  }
}

function updateFile() {
  updateTitle()
  if (fileSelector.value?.files) {
    const file = fileSelector.value.files[0]
    image.value = URL.createObjectURL(file)
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
            class="input border-none"
            @change="updateFile"
          />
        </div>
        <img
          v-if="image"
          :src="image"
          class="max-h-90vh max-w-90vw lg:max-w-60vw"
          alt="your image"
        />
        <div class="flex flex-col gap-4">
          <label for="title">Title:</label>
          <input
            id="title"
            v-model="title"
            type="text"
            name="title"
            :multiple="false"
            class="input border-solid"
            @change="updateTitle"
          />
        </div>
        <div>
          <button
            :disabled="uploadDisabled"
            class="button disabled:bg-red"
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
