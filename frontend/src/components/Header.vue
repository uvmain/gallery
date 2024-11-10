<script setup lang="ts">
import { useStorage } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getRandomSlug } from '../composables/getRandomSlug'

const props = defineProps({
  bg: { type: String, required: false },
  showAdd: { type: Boolean, default: false },
  showEdit: { type: Boolean, default: false },
})

const emit = defineEmits(['add', 'edit'])

const isModalOpened = ref(false)

function openModal() {
  isModalOpened.value = true
}

function closeModal() {
  isModalOpened.value = false
}

const bgColour = computed(() => {
  if (props.bg === '200')
    return 'bg-gray-200'
  else if (props.bg === '300')
    return 'bg-gray-300'
  else return 'bg-gray-100'
})

const router = useRouter()

const userLoginState = useStorage('login-state', false)

const iconColour = computed(() => {
  return userLoginState.value ? 'text-gray-600' : 'text-gray-400'
})

function enableEdit() {
  if (userLoginState.value) {
    emit('edit')
  }
  else {
    console.warn('Unable to enter edit mode, please log in')
  }
}

async function navigateHome() {
  router.push('/')
}

function navigateAlbums() {
  router.push('/albums')
}

async function navigateRandom() {
  const result = await getRandomSlug(1)
  const slug = result ? result[0] : ''
  router.push(`/${slug}`)
}

async function navigateUpload() {
  if (userLoginState.value)
    router.push('/upload')
}
</script>

<template>
  <div class="h-18 px-6" :class="bgColour">
    <header class="mx-auto flex justify-between lg:max-w-8/10 lg:p-6">
      <div class="flex gap-4">
        <div
          class="p-2 text-2xl text-gray-700 hover:cursor-pointer"
          @click="navigateHome"
        >
          home
        </div>
        <div
          class="p-2 text-2xl text-gray-700 hover:cursor-pointer"
          @click="navigateAlbums"
        >
          albums
        </div>
        <div
          class="p-2 text-2xl text-gray-700 hover:cursor-pointer"
          @click="navigateRandom"
        >
          random
        </div>
      </div>
      <div class="flex gap-4">
        <div class="p-2 hover:cursor-pointer" @click="navigateUpload">
          <icon-tabler-upload class="text-2xl" :class="iconColour" />
        </div>
        <slot />
        <div v-if="showAdd" class="p-2 hover:cursor-pointer" @click="emit('add')">
          <icon-tabler-library-plus class="text-2xl" :class="iconColour" />
        </div>
        <div v-if="showEdit" class="p-2 hover:cursor-pointer" @click="enableEdit">
          <icon-tabler-edit class="text-2xl" :class="iconColour" />
        </div>
        <div class="p-2 hover:cursor-pointer" @click="openModal">
          <icon-tabler-user class="text-2xl text-gray-600" />
        </div>
      </div>
      <LoginModal :is-open="isModalOpened" @modal-close="closeModal" @submit="navigateHome" />
    </header>
  </div>
</template>
