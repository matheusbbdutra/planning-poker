<template>
  <div class="bg-white p-8 rounded-lg shadow-md max-w-md mx-auto">

    <h2 class="text-2xl font-bold mb-2 text-center text-gray-800">
      Entrar na Sala
    </h2>
    <form @submit.prevent="handleJoinRoom">
      <div class="p-1">
        <label for="room-id" class="mb-2 text-sm font-medium text-gray-700">
          ID da Sala
        </label>

        <input
            id="room-id"
            v-model="roomId"
            type="text"
            placeholder="Informe o ID da sala"
            class="w-full p-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
        />
      </div>

      <div class="p-1">
        <label for="user-name" class="mb-2 text-sm font-medium text-gray-700">
          Seu nome
        </label>

        <input
            id="user-name"
            v-model="userName"
            type="text"
            placeholder="Digite seu nome para participar"
            class="w-full p-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
        />
      </div>

      <button
          type="submit"
          :disabled="!isFormValid"
          class="w-full mt-6 bg-green-600 text-white font-bold py-3 px-4 rounded-md hover:bg-green-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
      >
        Entrar e Votar
      </button>
    </form>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

const userName = ref('');
const roomId = ref('');

const isFormValid = computed(() => {
  return userName.value.trim() !== '' && roomId.value.trim() !== '';
});

const emit = defineEmits(['join-room']);

const handleJoinRoom = () => {
  if (!isFormValid.value) return;

  emit('join-room', {roomId: roomId.value.trim(), userName: userName.value.trim()});
};
</script>