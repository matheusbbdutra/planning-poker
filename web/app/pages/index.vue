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

const handleNewSession = async (sessionName: string, userName: string) => {
  const { data, error } = await useApi('/room', {
      method: 'POST',
      body: { sessionName, userName },
  });

  if (error.value) {
    console.error("Erro ao criar a sala:", error.value);
    return;
  }

  const room = data.value.room;
  const roomId = room?.id ?? room?.ID;
  const userId = room?.participants?.[0]?.id ?? room?.Participants?.[0]?.ID;

  if (!roomId || !userId) {
    console.error("Resposta inesperada ao criar sala:", room);
    return;
  }

  sessionStorage.setItem("userId", userId)
  router.push(`/room/${room.id}`);
};

const handleJoinRoom = (userName: string) => {
  console.log(`O usuário "${userName}" quer entrar na sala.`);
};
</script>