<script setup lang="ts">
import dayjs from 'dayjs'
import { onMounted, ref } from 'vue'
import { getServerUrl } from '../composables/getServerBaseUrl'

const serverBaseUrl = ref()
const username = ref('')
const password = ref('')
const loginStatus = ref(false)
const patchStatus = ref()

async function login() {
  serverBaseUrl.value = await getServerUrl()

  const formData = new FormData()
  formData.append('name', username.value)
  formData.append('password', password.value)

  const response = await fetch(`${serverBaseUrl.value}/api/login`, {
    body: formData,
    method: 'POST',
  })
  loginStatus.value = (response.status === 200)
}

async function logout() {
  serverBaseUrl.value = await getServerUrl()

  const response = await fetch(`${serverBaseUrl.value}/api/logout`, {
    method: 'GET',
    credentials: 'include',
  })
  loginStatus.value = (response.status !== 401)
}

async function checkIfLoggedIn() {
  try {
    serverBaseUrl.value = await getServerUrl()
    const response = await fetch(`${serverBaseUrl.value}/api/check-session`, {
      method: 'GET',
      credentials: 'include',
    })
    loginStatus.value = response.ok
  }
  catch {
    loginStatus.value = false
  }
}

async function patch() {
  serverBaseUrl.value = await getServerUrl()

  const patchData = {
    title: 'Soyuz Panel',
    // cameraMake: 'Nokia',
    // cameraModel: '7',
    dateTaken: dayjs('05-02-2017', 'MM-DD-YYYY'),
  }

  const slug = '1729943980387385100'

  const response = await fetch(`${serverBaseUrl.value}/api/metadata/${slug}`, {
    body: JSON.stringify(patchData),
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
  })
  patchStatus.value = response.status
}

onMounted(() => {
  checkIfLoggedIn()
})
</script>

<template>
  <Header />
  <div class="w-300 flex flex-col gap-4 p-6">
    <form>
      <label for="username">Username:</label>
      <input id="username" v-model="username" type="text" name="username" autocomplete="username"><br><br>
      <label for="password">Password:</label>
      <input id="password" v-model="password" type="password" name="password" autocomplete="current-password"><br><br>
    </form>
    <button @click="login">
      Login
    </button>
    <br>
    <button @click="checkIfLoggedIn">
      CheckLoginStatus
    </button>
    <br>
    <button @click="logout">
      Logout
    </button>
    <br>
    <p>
      login status: {{ loginStatus }}
    </p>
    <button @click="patch">
      Patch
    </button>
    {{ patchStatus }}
  </div>
</template>
