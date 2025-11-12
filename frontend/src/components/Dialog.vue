<script setup lang="ts">
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
  <dialog ref="dialog" class="modal p-2 md:p-4" @keydown.escape="dialog?.close()">
    <div :class="unpadded ? '' : 'p-4'">
      <slot />
    </div>
    <slot name="successOkButton">
      <button v-if="closeButton" aria-label="close" class="button" @click="dialog?.close()">
        Close
      </button>
    </slot>
  </dialog>
</template>
