<script setup lang="ts">
import { useRouter } from 'vue-router'
import { getThumbnailPath } from '../composables/logic'

const router = useRouter()

defineProps<{
  slug: string
}>()

function navigateToSlug(slug: string) {
  const slugPath = `/${slug}`
  router.push(slugPath)
}
</script>

<template>
  <div>
    <div class="absolute right-1 top-1 hover:cursor-pointer" @click="null">
      <!-- <icon-tabler-trash-x class="text-xl hover:size-115% text-white" /> -->
    </div>
    <img
      :src="getThumbnailPath(slug)"
      :alt="slug"
      onerror="this.onerror=null;this.src='/default-image.jpg';"
      class="h-40 w-80 cursor-pointer border-2 border-white border-solid object-cover dark:border-neutral-500"
      @click="navigateToSlug(slug)"
    />
  </div>
</template>

<style>
.grad {
  background: radial-gradient(circle at bottom left, rgb(0, 0, 0, 0.5) 0, rgba(255, 255, 255, 0) 60%);
}
</style>
