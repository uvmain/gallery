<script setup lang="ts">
import { onClickOutside, useSessionStorage } from '@vueuse/core'
import { backendFetchRequest } from '../composables/fetchFromBackend'

defineProps({
  isOpen: Boolean,
})

const emits = defineEmits(['modalClose'])

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
  emits('modalClose')
}

function cancel() {
  emits('modalClose')
}

async function logout() {
  const response = await backendFetchRequest('logout', {
    method: 'GET',
    credentials: 'include',
  })
  isLoggedIn.value = (response.status !== 401)
  userLoginState.value = (response.status !== 401)
  emits('modalClose')
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

onClickOutside(target, () => emits('modalClose'))
</script>

<template>
  <div v-if="isOpen" class="text size-full left-0 top-0 fixed z-999 backdrop-blur-xl">
    <div v-if="!isLoggedIn" @keydown.escape="cancel">
      <div ref="target" class="modal mx-auto mb-auto mt-150px px-30px pb-30px pt-20px w-300px">
        <div class="p-6 flex flex-col gap-4 w-300">
          <form id="login" class="flex flex-col gap-2">
            <div class="flex flex-row gap-2 items-center">
              <label for="username">Username:</label>
              <input id="username" v-model="username" type="text" name="username" autocomplete="username">
            </div>
            <div class="flex flex-row gap-2 items-center">
              <label for="password">Password:</label>
              <input id="password" v-model="password" type="password" name="password" autocomplete="current-password" @keydown.enter="login">
            </div>
          </form>
        </div>
        <div class="flex gap-4 justify-center">
          <button aria-label="cancel" class="button" @click="cancel">
            Cancel
          </button>
          <button aria-label="login" class="button" @click="login">
            Login
          </button>
        </div>
      </div>
    </div>
    <div v-else @keydown.escape="cancel">
      <div class="modal mx-auto mb-auto mt-150px px-30px pb-30px pt-20px rounded-sm w-300px">
        <div class="mb-2 py-4 text-center">
          You are logged in.
        </div>
        <div class="flex gap-4 justify-center">
          <button aria-label="cancel" class="button" @click="cancel">
            Cancel
          </button>
          <button aria-label="logout" class="button" @click="logout">
            Logout
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
