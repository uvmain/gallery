<script setup lang="ts">
defineProps({
  tooltipText: { type: String, required: true },
  content: { type: String, required: false },
})
</script>

<template>
  <div class="flex items-center space-x-3">
    <div class="group">
      <slot />
      <span class="tooltip">
        {{ tooltipText }}
      </span>
    </div>
    <span v-if="content">
      {{ content }}
    </span>
  </div>
</template>
