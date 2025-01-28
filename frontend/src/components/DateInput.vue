<script setup lang="ts">
defineProps({
  modelValue: {
    type: String,
    required: true,
  },
  displayValue: {
    type: String,
    required: false,
  },
  inEditingMode: {
    type: Boolean,
    default: false,
  },
})

const emits = defineEmits(['update:modelValue'])

function handleInput(event: any) {
  emits('update:modelValue', event.target.value)
}
</script>

<template>
  <div class="">
    <input
      v-show="inEditingMode"
      :placeholder="modelValue"
      :value="modelValue"
      type="datetime-local"
      class="border-1 border-white border-dashed p-1 outline-none standard"
      @input="handleInput"
    >
    <span
      v-show="!inEditingMode"
      class="py-2"
    >
      {{ displayValue }}
    </span>
  </div>
</template>
