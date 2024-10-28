<script setup lang="ts">
function addAlbum() {
  console.log('adding album')
}
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <Header @add="addAlbum" />
    <div id="main" class="flex flex-row justify-center gap-8 p-6">
      Add
    </div>
  </div>
</template>
