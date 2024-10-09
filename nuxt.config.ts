export default defineNuxtConfig({
  devtools: { enabled: true },
  compatibilityDate: '2024-10-09',
  modules: [
    "@unocss/nuxt",
    '@nuxt/eslint',
    "@nuxtjs/google-fonts",
    "nuxt-auth-utils",
    "@nuxt/icon"
  ],
  devServer: {
    port: 3000,
  },
  googleFonts: {
    families: {
      Poppins: true,
    }
  },
  imports: {
    dirs: [
      './types'
    ]
  }
})