<script setup lang="ts">
import { useRouter } from 'vue-router'

const props = defineProps({
  tag: { type: String, required: true },
  inEditingMode: { type: Boolean, default: false },
})

const router = useRouter()

function navigateToTag() {
  const tagPath = `/tags/${props.tag}`
  router.push(tagPath)
}
</script>

<template>
  <div class="text-nowrap tag" @click="navigateToTag">
    {{ tag }}
  </div>
</template>
