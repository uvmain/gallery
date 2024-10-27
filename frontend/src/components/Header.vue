<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { getRandomSlug } from '../composables/getRandomSlug'

const props = defineProps({
  bg: { type: String, required: true },
})

const bgColour = computed(() => {
  if (props.bg === '200')
    return 'bg-gray-200'
  else if (props.bg === '300')
    return 'bg-gray-300'
  else return 'bg-gray-100'
})

const router = useRouter()

async function navigateHome() {
  router.push('/')
}

async function navigateRandom() {
  const result = await getRandomSlug(1)
  const slug = result ? result[0] : ''
  router.push(`/${slug}`)
}
</script>

<template>
  <div class="h-18 px-6" :class="bgColour">
    <header class="mx-auto flex justify-between lg:max-w-8/10 lg:p-6">
      <div class="flex gap-4">
        <div
          class="p-2 text-xl text-gray-700 font-semibold hover:cursor-pointer"
          @click="navigateHome"
        >
          home
        </div>
        <div
          class="p-2 text-xl text-gray-700 font-semibold hover:cursor-pointer"
          @click="navigateHome"
        >
          albums
        </div>
        <div
          class="p-2 text-xl text-gray-700 font-semibold hover:cursor-pointer"
          @click="navigateRandom"
        >
          random
        </div>
      </div>
      <div class="flex gap-4">
        <div
          class="p-2 text-xl text-gray-700 font-semibold hover:cursor-pointer"
          @click="navigateHome"
        >
          <icon-tabler-edit class="text-3xl text-gray-600" />
        </div>
        <div
          class="p-2 text-xl text-gray-700 font-semibold hover:cursor-pointer"
          @click="navigateRandom"
        >
          <icon-tabler-user class="text-3xl text-gray-600" />
        </div>
      </div>
    </header>
  </div>
</template>
