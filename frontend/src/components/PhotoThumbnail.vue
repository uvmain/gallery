<script setup lang="ts">
import { getThumbnailPath } from '../composables/logic'

const props = defineProps({
  slug: { type: String, required: true },
  editMode: { type: Boolean, default: false },
  deleteMode: { type: Boolean, default: false },
  square: { type: Boolean, default: false },
  large: { type: Boolean, default: false },
})

const emits = defineEmits(['deleteImage', 'editImage'])

const router = useRouter()

const confirm = ref(false)

const sizeClass = computed(() => {
  if (props.square) {
    if (props.large) {
      return 'lg:size-80'
    }
    return 'lg:size-60'
  }
  else {
    if (props.large) {
      return 'lg:h-50 lg:w-100'
    }
    return 'lg:h-40 lg:w-80'
  }
})

function navigateToSlug(slug: string) {
  const slugPath = `/${slug}`
  router.push(slugPath)
}
</script>

<template>
  <div class="relative">
    <div v-if="deleteMode" class="grad absolute right-0 top-0 size-20 hover:cursor-pointer">
      <icon-tabler-trash-x v-if="!confirm" class="group absolute right-1 top-1 text-xl text-white hover:text-red" @click="confirm = true" />
      <icon-tabler-circle-dashed-check v-else class="group absolute right-1 top-1 text-xl text-white hover:text-green" @click="emits('deleteImage')" />
    </div>
    <div v-if="editMode" class="grad absolute right-0 top-0 size-20 hover:cursor-pointer">
      <icon-tabler-camera-rotate class="group absolute right-1 top-1 text-xl text-white hover:text-green" @click="emits('editImage')" />
    </div>
    <img
      :src="getThumbnailPath(slug)"
      :alt="slug"
      onerror="this.onerror=null;this.src='/default-image.jpg';"
      class="h-full w-full cursor-pointer border-2 border-white border-solid object-cover dark:border-neutral-500"
      :class="sizeClass"
      @click="navigateToSlug(slug)"
    />
  </div>
</template>

<style>
.grad {
  background: radial-gradient(circle at top right, rgb(0, 0, 0, 0.5) 0, rgba(255, 255, 255, 0) 60%);
}
</style>
