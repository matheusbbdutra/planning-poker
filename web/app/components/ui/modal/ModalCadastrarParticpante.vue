<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
    <div class="bg-white p-6 rounded-lg w-full max-w-md">
      <h3 class="text-lg font-bold mb-4">Identifique-se</h3>
      <label class="block text mb-2">Seu nome</label>
      <input v-model="name" class="w-full p-3 border rounded mb-4" placeholder="Digite seu nome" />
      <div class="flex justify-end gap-2">
        <button @click="onCancel" class="px-4 py-2 rounded bg-gray-200">Cancelar</button>
        <button :disabled="!name.trim()" @click="onRegister" class="px-4 py-2 rounded bg-blue-600 text-white" >
          Confirmar
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
const emit = defineEmits<{
  (e: 'registered', name: string): void;
  (e: 'close'): void;
}>();

const name = ref('');

const onRegister = () => {
  emit('registered', name.value.trim() );
};

const onCancel = () => {
  emit('close');
};
</script>