<script setup lang="ts">
const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'add-task', task: string): void;
}>();
const input = ref('');
const error = ref<string | null>(null);

const handleAddTarefa = () => {
  if (input.value.trim()) {
    emit('add-task', input.value.trim());
  }

  emit('close');
};

const emitClose = () => {
  emit('close');
};
</script>

<template>
  <div class="fixed inset-0 flex items-center justify-center z-50">
    <!-- Backdrop levemente escuro -->
    <div
      class="absolute inset-0 bg-black/20 transition-opacity"
      @click="emitClose"
      aria-hidden="true"
    ></div>

    <!-- Modal acima do backdrop -->
    <div class="relative bg-white rounded-lg shadow-lg w-full max-w-md p-6">
      <h2 class="text-2xl font-bold mb-4">Adicionar Tarefa</h2>
      <form @submit.prevent="handleAddTarefa">
        <div class="mb-4">
          <label for="titulo" class="block text-sm font-medium text-gray-700 mb-1">Título</label>
          <input  v-model="input" type="text" id="titulo" class="w-full border border-gray-300 rounded-md p-2 focus:outline-none focus:ring-2 focus:ring-blue-500" />
        </div>
        <div class="flex justify-end space-x-2">
          <button type="button" @click="emitClose" class="px-4 py-2 bg-gray-200 rounded-md hover:bg-gray-300">Cancelar</button>
          <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">Salvar</button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>

</style>