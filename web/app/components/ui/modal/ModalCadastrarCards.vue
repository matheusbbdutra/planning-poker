<template>
  <!-- ...existing code... -->
  <div class="fixed inset-0  bg-opacity-50 flex items-center justify-center z-50" @keydown.esc="emitClose" tabindex="-1">
     <div
      class="absolute inset-0 bg-black/20 transition-opacity"
      @click="emitClose"
      aria-hidden="true"
    ></div>

    <div class="relative bg-white rounded-lg shadow-lg w-full max-w-md p-6" role="dialog" aria-modal="true" aria-labelledby="modal-title">
      <h2 id="modal-title" class="text-2xl font-bold mb-4">Editar cartas</h2>
      <form @submit.prevent="handleCreateCard">
        <div class="mb-4">
          <label for="cartasInput" class="block text-sm font-medium text-gray-700 mb-1">Preencha as cartas separadas por vírgula (ex: 1,2,3,5,8,13,☕,?):</label>
          <input
            v-model="input"
            type="text"
            id="cartasInput"
            class="w-full border border-gray-300 rounded-md p-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            aria-describedby="cartasHelp"
          />
          <p id="cartasHelp" class="text-xs text-gray-500 mt-1">Valores não vazios serão considerados. Símbolos também são permitidos.</p>
          <p v-if="error" class="text-sm text-red-600 mt-2">{{ error }}</p>
        </div>
        <div class="flex justify-end space-x-2">
          <button type="button" @click="emitClose" class="px-4 py-2 bg-gray-200 rounded-md hover:bg-gray-300">Cancelar</button>
          <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">Salvar</button>
        </div>
      </form>
    </div>
  </div>
  <!-- ...existing code... -->
</template>

<script lang="ts" setup>
// ...existing code...
import { ref } from 'vue';

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'save', cards: Array<string | number>): void;
}>();

const input = ref('');
const error = ref<string | null>(null);

const parseCards = (raw: string): Array<string | number> => {
  return raw
    .split(',')
    .map(s => s.trim())
    .filter(s => s.length > 0)
    .map(s => {
      // tenta converter para número quando aplicável
      const n = Number(s);
      return Number.isFinite(n) ? n : s;
    });
};

const handleCreateCard = () => {
  error.value = null;
  const parsed = parseCards(input.value);
  if (parsed.length === 0) {
    error.value = 'Informe ao menos uma carta válida.';
    return;
  }

  // Emite as cartas para o pai (para atualizar estado/global)
  emit('save', parsed);

  // Fecha o modal
  emit('close');
  input.value = '';
};

const emitClose = () => {
  emit('close');
};
// ...existing code...
</script>

<style>
/* ...existing code... */
</style>