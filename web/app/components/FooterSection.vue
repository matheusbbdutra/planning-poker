<template>
  <footer class="bg-gray-900">
    <div class="container mx-auto py-6 text-center">
      <p class="text-sm text-gray-400">
        &copy; {{ currentYear }} {{ authorName }}. All rights reserved.
      </p>
    </div>
  </footer>
</template>

<script setup lang="ts">

import { computed } from 'vue';

defineProps({
  authorName: {
    type: String,
    required: true, 
  }
});
const currentYear = computed(() => new Date().getFullYear());
</script>