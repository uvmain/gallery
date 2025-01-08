<script setup>
import { onClickOutside, useSessionStorage } from '@vueuse/core'
import { onBeforeMount, ref } from 'vue'

import { backendFetchRequest } from '../composables/fetchFromBackend'

defineProps({
  isOpen: Boolean,
})

const emit = defineEmits(['modalClose'])
const username = ref('')
const password = ref('')
const isLoggedIn = ref(false)
const target = ref(null)

const userLoginState = useSessionStorage('login-state', isLoggedIn.value)

async function login() {
  const formData = new FormData()
  formData.append('username', username.value)
  formData.append('password', password.value)

  const response = await backendFetchRequest('login', {
    body: formData,
    method: 'POST',
  })
  isLoggedIn.value = (response.status !== 401)
  userLoginState.value = (response.status !== 401)
  emit('modalClose')
}

function cancel() {
  emit('modalClose')
}

async function logout() {
  const response = await backendFetchRequest('logout', {
    method: 'GET',
    credentials: 'include',
  })
  isLoggedIn.value = (response.status !== 401)
  userLoginState.value = (response.status !== 401)
  emit('modalClose')
}

async function checkIfLoggedIn() {
  try {
    const response = await backendFetchRequest('check-session', {
      method: 'GET',
      credentials: 'include',
    })
    if (response.status === 401) {
      isLoggedIn.value = false
      userLoginState.value = false
    }
    isLoggedIn.value = response.ok
    userLoginState.value = response.ok
  }
  catch {
    isLoggedIn.value = false
    userLoginState.value = false
  }
}

onBeforeMount(() => {
  checkIfLoggedIn()
})

onClickOutside(target, () => emit('modalClose'))
</script>

<template>
  <div v-if="isOpen" class="fixed left-0 top-0 z-999 size-full backdrop-blur-xl">
    <div v-if="!isLoggedIn">
      <div ref="target" class="mx-auto mb-auto mt-150px w-300px rounded-lg bg-neutral-200 px-30px pb-30px pt-20px text-neutral-700 dark:bg-neutral-700 dark:text-neutral-100">
        <div class="w-300 flex flex-col gap-4 p-6">
          <form class="flex flex-col gap-2">
            <div class="flex flex-row items-center gap-2">
              <label for="username">Username:</label>
              <input id="username" v-model="username" type="text" name="username" autocomplete="username">
            </div>
            <div class="flex flex-row items-center gap-2">
              <label for="password">Password:</label>
              <input id="password" v-model="password" type="password" name="password" autocomplete="current-password" @keydown.enter="login">
            </div>
          </form>
        </div>
        <div class="flex justify-center gap-4">
          <button aria-label="cancel" class="px-4 py-2" @click="cancel">
            Cancel
          </button>
          <button aria-label="login" class="px-4 py-2" @click="login">
            Login
          </button>
        </div>
      </div>
    </div>
    <div v-else>
      <div class="mx-auto mb-auto mt-150px w-300px rounded-lg px-30px pb-30px pt-20px standard">
        <div class="mb-2 py-4 text-center">
          You are logged in.
        </div>
        <div class="flex justify-center gap-4">
          <button aria-label="cancel" class="px-4 py-2" @click="cancel" @keydown.escape="cancel">
            Cancel
          </button>
          <button aria-label="logout" class="px-4 py-2" @click="logout">
            Logout
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
