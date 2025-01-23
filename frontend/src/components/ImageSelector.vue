<script setup lang="ts">
import { useSessionStorage } from '@vueuse/core'
import { computed, onMounted, ref } from 'vue'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import { getThumbnailPath } from '../composables/logic'

const props = defineProps({
  singleSelect: { type: Boolean, default: false },
  definedSlugs: { type: Array, required: false },
})

const emits = defineEmits(['closeModal'])

const slugs = ref()

const selectedSlugs = useSessionStorage<string[]>('selected-slugs', [])

async function getSlugs() {
  if (props.definedSlugs?.length) {
    slugs.value = props.definedSlugs as string[]
  }
  else {
    const response = await backendFetchRequest('slugs')
    const jsonData = await response.json() as string[]
    slugs.value = jsonData
  }
}

function toggleSelected(slug: string) {
  const index = selectedSlugs.value.indexOf(slug)
  if (index > -1) {
    selectedSlugs.value.splice(index, 1)
  }
  else {
    selectedSlugs.value.push(slug)
    if (props.singleSelect === true) {
      emits('closeModal', slug)
    }
  }
}

const slugsToUse = computed(() => {
  return props.definedSlugs?.length ? props.definedSlugs as string[] : slugs.value
})

onMounted(() => {
  getSlugs()
})
</script>

<template>
  <div class="mt-4 flex flex-wrap gap-x-2 gap-y-1">
    <div v-for="(slug, index) in slugsToUse" :key="index" class="relative">
      <icon-tabler-circle-check-filled v-if="selectedSlugs.includes(slug)" class="absolute right-1 top-1 rounded-2xl bg-white text-2xl text-green" />
      <img :src="getThumbnailPath(slug)" :alt="slug" class="size-40 cursor-pointer" @click="toggleSelected(slug)">
    </div>
  </div>
</template>
