<script setup lang="ts">
import { useDark, useSessionStorage, useToggle } from '@vueuse/core'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getRandomSlug } from '../composables/getRandomSlug'

defineProps({
  showAdd: { type: Boolean, default: false },
  showEdit: { type: Boolean, default: false },
})

const emit = defineEmits(['add', 'edit'])

const isModalOpened = ref(false)
const isDark = useDark()
const toggleDark = useToggle(isDark)

function openModal() {
  isModalOpened.value = true
}

function closeModal() {
  isModalOpened.value = false
}

const router = useRouter()

const userLoginState = useSessionStorage('login-state', false)

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
  <div class="h-18 px-6 standard">
    <header class="mx-auto flex justify-between lg:max-w-8/10 lg:p-6">
      <div class="flex gap-4">
        <div
          class="p-2 text-2xl hover:cursor-pointer"
          @click="navigateHome"
        >
          home
        </div>
        <div
          class="p-2 text-2xl hover:cursor-pointer"
          @click="navigateAlbums"
        >
          albums
        </div>
        <div
          class="p-2 text-2xl hover:cursor-pointer"
          @click="navigateRandom"
        >
          random
        </div>
      </div>
      <div class="flex gap-4">
        <div class="p-2 hover:cursor-pointer" @click="navigateUpload">
          <icon-tabler-upload class="text-2xl" />
        </div>
        <slot />
        <div v-if="showAdd" class="p-2 hover:cursor-pointer" @click="emit('add')">
          <icon-tabler-library-plus class="text-2xl" />
        </div>
        <div v-if="showEdit" class="p-2 hover:cursor-pointer" @click="enableEdit">
          <icon-tabler-edit class="text-2xl" />
        </div>
        <div v-if="isDark" class="p-2 hover:cursor-pointer" @click="toggleDark()">
          <icon-tabler-moon-stars class="text-2xl" />
        </div>
        <div v-else class="p-2 hover:cursor-pointer" @click="toggleDark()">
          <icon-tabler-sun class="text-2xl" />
        </div>
        <div class="p-2 hover:cursor-pointer" @click="openModal">
          <icon-tabler-user class="text-2xl" />
        </div>
      </div>
      <LoginModal :is-open="isModalOpened" @modal-close="closeModal" @submit="navigateHome" />
    </header>
  </div>
</template>
