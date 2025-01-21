<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const props = defineProps({
  imageSlug: { type: String, required: true },
  inEditingMode: { type: Boolean, default: false },
})

const tags = ref<string[]>([])

async function getTags() {
  tags.value = []
  const response = await backendFetchRequest(`tags/${props.imageSlug}`)
  tags.value = await response.json() || []
}

async function addTag() {
  return null
}

watch(
  () => props.imageSlug,
  () => {
    getTags()
  },
)

onMounted(() => {
  getTags()
})
</script>

<template>
  <div class="mt-4 max-w-120 flex flex-col gap-4 border-1 border-gray-400 rounded-sm border-solid p-4 dark:border-gray-600">
    <div class="text-left text-lg">
      Tags
    </div>
    <div class="flex flex-wrap gap-4">
      <div v-for="(tag, index) in tags" :key="index">
        <Tag :tag="tag" />
      </div>
      <div v-if="inEditingMode">
        <div class="hover:cursor-pointer" @click="addTag()">
          <icon-tabler-square-plus class="text-2xl" />
        </div>
      </div>
    </div>
  </div>
</template>
