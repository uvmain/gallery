<script setup lang="ts">
import type Dialog from '../components/Dialog.vue'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const props = defineProps({
  imageSlug: { type: String, required: true },
  inEditingMode: { type: Boolean, default: false },
})

defineExpose({ getTags })

const tags = ref<string[]>([])
const newTagString = ref<string>()
const addDialog = ref<typeof Dialog>()

async function getTags() {
  tags.value = []
  const tagsRequest = await backendFetchRequest(`tags/${props.imageSlug}`)
  tags.value = await tagsRequest.json() || []
}

function showAddDialog() {
  addDialog.value?.show()
}

function hideAddDialog() {
  addDialog.value?.hide()
}

async function confirmAddTag() {
  const newTag = {
    Tag: newTagString.value,
    ImageSlugs: [props.imageSlug],
  }
  const options = {
    body: JSON.stringify(newTag),
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest('tags', options)

  hideAddDialog()

  if (response.status === 200) {
    await getTags()
  }
}

async function deleteTag(tag: string) {
  const newTag = {
    Tag: tag,
    ImageSlug: props.imageSlug,
  }
  const options = {
    body: JSON.stringify(newTag),
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest('tags', options)

  hideAddDialog()

  if (response.status === 200) {
    await getTags()
  }
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
  <div class="mt-4 p-4 border-1 border-gray-400 rounded-sm border-solid flex flex-col gap-4 max-w-120 dark:border-gray-600">
    <div class="flex flex-wrap gap-4">
      <div class="text-lg text-left">
        Tags
      </div>
      <div v-for="(tag, index) in tags" :key="index">
        <Tag :tag="tag" :in-editing-mode="inEditingMode" @delete="deleteTag(tag)" />
      </div>
      <div v-if="inEditingMode">
        <div class="hover:cursor-pointer" @click="showAddDialog()">
          <icon-tabler-square-plus class="text-2xl" />
        </div>
      </div>
    </div>
    <Dialog ref="addDialog" :close-button="false" class="modal">
      <div class="p-6">
        <div class="flex flex-col gap-6 items-center">
          <icon-tabler-library-plus class="text-4xl text-green" />
          <div class="flex flex-row gap-2 items-center">
            <label for="newtag">New Tag:</label>
            <input id="newtag" v-model="newTagString" type="text" name="newtag" @keydown.enter="confirmAddTag()">
          </div>
          <div class="flex gap-4 justify-center">
            <button aria-label="cancel" class="button" @click="hideAddDialog()">
              Cancel
            </button>
            <button aria-label="delete" class="button" @click="confirmAddTag()">
              Add
            </button>
          </div>
        </div>
      </div>
    </Dialog>
  </div>
</template>
