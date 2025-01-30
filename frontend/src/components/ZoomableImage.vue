<script setup lang="ts">
defineProps({
  imageSource: { type: String, required: true },
  width: { type: Number, required: false, default: 0 },
  height: { type: Number, required: false, default: 0 },
})

const image = ref<HTMLImageElement>()
const zoomed = ref(false)

const zoomedClass = computed(() => {
  return zoomed.value
    ? 'fixed inset-0 m-auto max-h-screen max-w-screen'
    : 'h-auto w-auto max-h-100vh max-w-95vw lg:max-h-80vh lg:max-w-70vw border-8 border-white border-solid box-border dark:border-neutral-500'
})

function toggleZoomed() {
  zoomed.value = !zoomed.value
}
</script>

<template>
  <div v-if="imageSource">
    <div v-if="zoomed" class="fixed inset-0 z-40 bg-black bg-opacity-75" @click="toggleZoomed()" />
    <img
      ref="image"
      :src="imageSource"
      :width="width"
      :height="height"
      class="z-50 cursor-pointer object-contain"
      :class="zoomedClass"
      @click="toggleZoomed()"
    />
  </div>
</template>
