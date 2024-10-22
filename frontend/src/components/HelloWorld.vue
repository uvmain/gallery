<script setup lang="ts">
defineProps({
  message: { type: String, required: false, default: 'Hello World' },
})
</script>

<template>
  <div class="h-screen flex items-center justify-center bg-dark">
    <h1 class="text-fuchsia-700">
      {{ message }}
    </h1>
  </div>
</template>
