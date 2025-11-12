<script setup lang="ts">
defineProps({
  modelValue: {
    type: String,
    required: true,
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
  <div class="w-full">
    <input
      v-show="inEditingMode"
      :placeholder="modelValue"
      :value="modelValue"
      type="text"
      class="w-full border-1 border-white border-dashed standard p-1"
      @input="handleInput"
    >
    <span
      v-show="!inEditingMode"
      class="py-2"
    >
      {{ modelValue }}
    </span>
  </div>
</template>
