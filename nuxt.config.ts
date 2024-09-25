export default defineNuxtConfig({
  devtools: { enabled: true },
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
  },
  runtimeConfig: {
    adminUser: 'admin', // can be overridden by NUXT_ADMIN_USER environment variable
    adminPassword: 'password', // can be overridden by NUXT_ADMIN_PASSWORD environment variable
    sessionMaxAge: 86400, // cookie session length in seconds, can be overridden by NUXT_SESSION_MAX_AGE environment variable
    databasePath: './database' // location of sqlite.db, can be overridden by NUXT_DATABASE_PATH environment variable
  },
});