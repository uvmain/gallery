<script setup>

useHead({
  htmlAttrs: {
    lang: 'en',
  },
  title: 'Gallery',
})

</script>

<template>
  <div id="app" class="bg-blue-gray-100 min-h-screen flex flex-col">
    <Header />
    <NuxtPage class="flex-1 text-dark text-base"/>
  </div>
</template>

<style>
html, body, #__nuxt {
  margin: 0;
  padding: 0;
  border: 0;
  font-family: 'Poppins', sans-serif;
}
</style>