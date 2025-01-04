<script lang="ts" setup>
import { ref } from 'vue'

defineProps({
  title: { type: String, required: false, default: null },
  closeButton: { type: Boolean, required: false, default: false },
  unpadded: { type: Boolean, required: false, default: false },
})

const dialog = ref<HTMLDialogElement>()

function show() {
  dialog.value?.showModal()
}

function hide() {
  dialog.value?.close()
}

defineExpose({ show, hide })
</script>

<template>
  <dialog ref="dialog" class="rounded-lg bg-transparent px-2 pb-4 md:pb-8">
    <div :class="unpadded ? '' : 'px-4 pt-4'">
      <slot />
    </div>
    <slot name="successOkButton">
      <button v-if="closeButton" class="mx-auto mt-4 block rounded px-6 py-3" @click="dialog?.close()">
        Close
      </button>
    </slot>
  </dialog>
</template>
