import type { Component } from 'vue'
import generatedRoutes from 'virtual:generated-pages'
import { ViteSSG } from 'vite-ssg'

import App from './App.vue'
import 'virtual:uno.css'

export const createApp = ViteSSG(App as Component, { routes: generatedRoutes, base: import.meta.env.BASE_URL })
