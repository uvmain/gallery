<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const props = defineProps({
  imageSlug: { type: String, required: true },
  inEditingMode: { type: Boolean, default: false },
})

defineExpose({ getTags })

const tags = ref<string[]>([])

async function getTags() {
  tags.value = []
  const tagsRequest = await backendFetchRequest(`tags/${props.imageSlug}`)
  const returnArray = await tagsRequest.json() || []

  const dimensionsRequest = await backendFetchRequest(`dimensions/${props.imageSlug}`)
  const dimensionsResponse = await dimensionsRequest.json()

  const albumsRequest = await backendFetchRequest(`links/image/${props.imageSlug}`)
  if (albumsRequest != null) {
    const albumSlugs: string[] = await albumsRequest.json() || []
    for (const albumSlug of albumSlugs) {
      const albumRequest: any = await backendFetchRequest(`albums/${albumSlug}`)
      const album = await albumRequest.json()
      returnArray.push(album.Name)
    }
  }

  if (dimensionsResponse.Panoramic === true) {
    returnArray.push('panoramic')
  }
  returnArray.push(dimensionsResponse.Orientation)
  tags.value = [...new Set(returnArray)] as string[]
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
