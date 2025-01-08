<script setup lang="ts">
import { useDark, useSessionStorage, useToggle } from '@vueuse/core'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import TooltipIcon from '../components/TooltipIcon.vue'
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
        <TooltipIcon v-if="userLoginState" tooltip-text="Upload" class="hover:cursor-pointer" @click="navigateUpload">
          <icon-tabler-upload class="text-2xl" />
        </TooltipIcon>
        <slot />
        <TooltipIcon v-if="showAdd && userLoginState" tooltip-text="Add" class="hover:cursor-pointer" @click="emit('add')">
          <icon-tabler-library-plus class="text-2xl" />
        </TooltipIcon>
        <TooltipIcon v-if="showEdit && userLoginState" tooltip-text="Edit Mode" class="hover:cursor-pointer" @click="enableEdit">
          <icon-tabler-edit class="text-2xl" />
        </TooltipIcon>
        <TooltipIcon :tooltip-text="isDark ? 'Light Mode' : 'Dark Mode'" class="hover:cursor-pointer" @click="toggleDark()">
          <icon-tabler-sun v-if="isDark" class="text-2xl" />
          <icon-tabler-moon-stars v-else class="text-2xl" />
        </TooltipIcon>
        <TooltipIcon :tooltip-text="userLoginState ? 'Log Out' : 'Log In'" class="hover:cursor-pointer" @click="openModal">
          <icon-tabler-user class="text-2xl" />
        </TooltipIcon>
      </div>
      <LoginModal :is-open="isModalOpened" @modal-close="closeModal" @submit="navigateHome" />
    </header>
  </div>
</template>
