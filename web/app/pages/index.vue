<template>
  <section class="bg-white min-h-screen flex">
    <div class="container mx-auto px-6 py-10">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
        <div>
          <h1 class="text-4xl lg:text-5xl font-bold text-gray-800 mb-4">
            Planning Poker Online e Descomplicado
          </h1>
          <p class="text-lg text-gray-600 mb-8">
            Crie uma sala, convide sua equipe e estime tarefas de forma rápida e colaborativa. Sem login, sem complicação.
          </p>
          <img src="~/assets/images/group.svg" alt="Ilustração de uma equipe colaborando" class="w-full h-4/12 rounded-lg">
        </div>
        <div>
          <div class="p-6">
            <CreateRoomCard @create-room="handleNewSession" />
          </div>
            <div class="p-6">
              <JoinRoomCard @roomId="handleJoinRoom" />
            </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import CreateRoomCard from "@/components/ui/cards/CreateRoomCard.vue";
import { useRouter } from '#vue-router';
import JoinRoomCard from "@/components/ui/cards/JoinRoomCard.vue"; 

onMounted(() => {
  document.title = "Planning Poker - Estime Tarefas com Sua Equipe";
});

const router = useRouter();

const handleNewSession = (sessionName: string, userName: string) => {
  fetch("/api/room/create", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ sessionName, userName }),
  })
    .then((response) => response.json())
    .then((data) => {
      const room = data.room;
      sessionStorage.setItem("roomId", room.id);
      sessionStorage.setItem("userId", room.participants[0].id);
      router.push(`/room/${room.id}`);
    })
    .catch((error) => {
      console.error("Erro ao criar a sala:", error);
    });

};

const handleJoinRoom = (userName: string) => {
  console.log(`O usuário "${userName}" quer entrar na sala.`);
};
</script>