import { createPinia } from 'pinia'
import routes from 'virtual:generated-pages'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'

import App from './App.vue'
import 'virtual:uno.css'

// eslint-disable-next-line ts/no-unsafe-argument
const app = createApp(App)

app.use(createPinia())

const router = createRouter({
  history: createWebHistory(),
  routes,
})

app.use(router)

app.mount('#app')
