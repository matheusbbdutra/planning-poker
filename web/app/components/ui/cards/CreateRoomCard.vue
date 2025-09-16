<template>
  <div class="bg-white p-8 rounded-lg shadow-md max-w-md mx-auto">

    <h2 class="text-2xl font-bold mb-6 text-center text-gray-800">
      Criar Nova Sala 🃏
    </h2>

    <form @submit.prevent="handleCreateRoom">
      <div class="p-1">
        <label for="name" class="block text-sm font-medium text-gray-700">
          Seu nome
        </label>
        <input
            id="name"
            v-model="userName"
            type="text"
            placeholder="Digite seu nome"
            class="w-full p-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
        />
      </div>
      <div class="p-1">
        <label for="session-name" class="block mb-2 text-sm font-medium text-gray-700">
          Nome da sessão
        </label>

        <input
            id="session-name"
            v-model="sessionName"
            type="text"
            placeholder="Ex: Sprint 25 - Novas Features"
            class="w-full p-3 border border-gray-300 rounded-md  focus:ring-blue-500 focus:border-blue-500 transition"
        />
      </div>

      <button
          type="submit"
          :disabled="!isFormValid"
          class="w-full mt-6 bg-blue-600 text-white font-bold py-3 px-4 rounded-md hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
      >
        Criar Sala
      </button>
    </form>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

const sessionName = ref('');
const userName = ref('');

const isFormValid = computed(() => {
  return sessionName.value.trim() !== '' && userName.value.trim() !== '';
});

const emit = defineEmits(['create-room']);

const handleCreateRoom = () => {
  if (!isFormValid.value) return; 
  emit('create-room', sessionName.value.trim(), userName.value.trim());
};
</script>
