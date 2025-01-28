<script setup lang="ts">
import type Dialog from './Dialog.vue'

const props = defineProps({
  tag: { type: String, required: true },
  inEditingMode: { type: Boolean, default: false },
})

const emits = defineEmits(['delete'])

const router = useRouter()
const deleteDialog = ref<typeof Dialog>()

function onClick() {
  if (!props.inEditingMode) {
    const tagPath = `/tags/${props.tag}`
    router.push(tagPath)
  }
  else {
    deleteDialog.value?.show()
  }
}

function hideDeleteDialog() {
  deleteDialog.value?.hide()
}
</script>

<template>
  <div>
    <div class="text-nowrap tag" :class="inEditingMode ? 'hover:bg-red hover:dark:bg-red' : ''" @click="onClick">
      {{ tag }}
    </div>
    <Dialog ref="deleteDialog" :close-button="false" class="modal">
      <div class="p-6">
        <div class="flex flex-col items-center gap-6">
          <p>
            Delete Tag: {{ tag }}
          </p>
          <div class="flex justify-center gap-4">
            <button aria-label="cancel" class="button" @click="hideDeleteDialog()">
              Cancel
            </button>
            <button aria-label="delete" class="button" @click="emits('delete', tag)">
              Delete
            </button>
          </div>
        </div>
      </div>
    </Dialog>
  </div>
</template>
