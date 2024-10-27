<script setup>
import { onClickOutside, useStorage } from '@vueuse/core'
import { onBeforeMount, ref } from 'vue'

import { getServerUrl } from '../composables/getServerBaseUrl'

defineProps({
  isOpen: Boolean,
})

const emit = defineEmits(['modalClose'])
const serverBaseUrl = ref()
const username = ref('')
const password = ref('')
const isLoggedIn = ref(false)
const target = ref(null)

const userLoginState = useStorage('login-state', isLoggedIn.value)

async function login() {
  serverBaseUrl.value = await getServerUrl()

  const formData = new FormData()
  formData.append('name', username.value)
  formData.append('password', password.value)

  const response = await fetch(`${serverBaseUrl.value}/api/login`, {
    body: formData,
    method: 'POST',
  })
  isLoggedIn.value = (response.status !== 401)
  userLoginState.value = (response.status !== 401)
}

async function logout() {
  serverBaseUrl.value = await getServerUrl()

  const response = await fetch(`${serverBaseUrl.value}/api/logout`, {
    method: 'GET',
    credentials: 'include',
  })
  isLoggedIn.value = (response.status !== 401)
  userLoginState.value = (response.status !== 401)
}

async function checkIfLoggedIn() {
  try {
    serverBaseUrl.value = await getServerUrl()
    const response = await fetch(`${serverBaseUrl.value}/api/check-session`, {
      method: 'GET',
      credentials: 'include',
    })
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
  <div v-if="isOpen" class="fixed left-0 top-0 z-999 size-full bg-dark">
    <div v-if="!isLoggedIn">
      <div ref="target" class="mx-auto mb-auto mt-150px w-300px rounded-lg bg-white px-30px pb-30px pt-20px shadow-dark shadow-xl">
        <div class="w-300 flex flex-col gap-4 p-6">
          <form class="flex flex-col gap-2">
            <div class="flex flex-row gap-2">
              <label for="username">Username:</label>
              <input id="username" v-model="username" type="text" name="username" autocomplete="username">
            </div>
            <div class="flex flex-row gap-2">
              <label for="password">Password:</label>
              <input id="password" v-model="password" type="password" name="password" autocomplete="current-password">
            </div>
          </form>
        </div>
        <div class="flex justify-center">
          <button @click="login" @click.stop="emit('modalClose')">
            Login
          </button>
        </div>
      </div>
    </div>
    <div v-else>
      You are logged in;
      <div class="flex justify-center">
        <button @click="logout" @click.stop="emit('modalClose')">
          Logout
        </button>
      </div>
    </div>
  </div>
</template>
