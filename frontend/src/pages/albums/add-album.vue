<script setup lang="ts">
import { ref } from 'vue'
import { backendFetchRequest } from '../../composables/fetchFromBackend'

const newAlbumResponse = ref()

async function addAlbum() {
  const newAlbum = {
    Name: 'Macro',
    CoverSlug: '1729943979792078600',
  }
  const options = {
    body: JSON.stringify(newAlbum),
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest('albums', options)

  newAlbumResponse.value = response.status

  // if (response.status === 200) {
  //   albums.value = await getAllAlbums()
  // }
}
</script>

<template>
  <div class="min-h-screen bg-gray-300">
    <Header :show-edit="true">
    </Header>
    <div @click="addAlbum">
      testing
    </div>
    <div v-if="newAlbumResponse">
      {{ newAlbumResponse }}
    </div>
  </div>
</template>
