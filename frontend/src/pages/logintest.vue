<script setup lang="ts">
import { ref } from 'vue'
import { getServerUrl } from '../composables/getServerBaseUrl'

const serverBaseUrl = ref()
const username = ref('')
const password = ref('')
const loginText = ref('')
const loginStatus = ref()

async function login() {
  serverBaseUrl.value = await getServerUrl()

  const formData = new FormData()
  formData.append('name', username.value)
  formData.append('password', password.value)

  const response = await fetch(`${serverBaseUrl.value}/api/login`, {
    body: formData,
    method: 'post',
  })
  loginStatus.value = response.status
  loginText.value = await response.text()
}
</script>

<template>
  <div class="w-300 flex flex-col gap-4 p-6">
    <label for="username">Username:</label>
    <input id="username" v-model="username" type="text" name="username"><br><br>
    <label for="password">Password:</label>
    <input id="password" v-model="password" type="password" name="password"><br><br>
    <button @click="login">
      Login
    </button>
  </div>
  <p>Status: {{ loginStatus }}</p>
  <p>Response: {{ loginText }}</p>
</template>
