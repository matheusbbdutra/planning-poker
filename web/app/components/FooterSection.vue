<template>
  <footer class="bg-slate-950 text-slate-200">
    <div class="mx-auto flex w-full max-w-6xl flex-col gap-10 px-6 py-12 md:flex-row md:items-start md:justify-between">
      <div class="space-y-4 md:max-w-sm">
        <h2 class="text-xl font-semibold text-white">Planning Poker</h2>
        <p class="text-sm text-slate-400">
          Experiência colaborativa para estimar tarefas em horas, com backend em Go, frontend Nuxt e Redis para tempo real.
        </p>
        <div class="flex flex-wrap gap-3 text-sm">
          <a
            v-if="props.projectLink"
            :href="props.projectLink"
            target="_blank"
            rel="noopener"
            class="inline-flex items-center rounded-full border border-slate-200 px-4 py-2 font-medium text-slate-200 transition hover:bg-slate-900"
          >
            Ver projeto
          </a>
        </div>
      </div>

      <div class="grid gap-8 text-sm sm:grid-cols-2">
        <div class="space-y-3">
          <h3 class="text-xs font-semibold uppercase tracking-wide text-slate-400">Navegação</h3>
          <ul class="space-y-2">
            <li>
              <NuxtLink to="/" class="transition hover:text-white">Início</NuxtLink>
            </li>
            <li>
              <a href="#como-funciona" class="transition hover:text-white">Como funciona</a>
            </li>
            <li>
              <a href="#documentacao" class="transition hover:text-white">Guia rápido</a>
            </li>
          </ul>
        </div>
        <div class="space-y-3">
          <h3 class="text-xs font-semibold uppercase tracking-wide text-slate-400">Recursos</h3>
          <ul class="space-y-2">
            <li>
              <a
                href="https://nuxt.com/docs/getting-started/introduction"
                target="_blank"
                rel="noopener"
                class="transition hover:text-white"
              >
                Documentação Nuxt
              </a>
            </li>
            <li v-if="props.contactEmail">
              <a
                :href="`mailto:${props.contactEmail}`"
                class="transition hover:text-white"
              >
                Fale conosco
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>

    <div class="border-t border-slate-900/60 py-4 text-center text-xs text-slate-500">
      © {{ currentYear }} {{ props.authorName }}. Todos os direitos reservados.
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed, withDefaults } from 'vue';

const props = withDefaults(defineProps<{
  authorName: string;
  projectLink?: string;
  contactEmail?: string;
}>(), {
  projectLink: '',
  contactEmail: '',
});

const currentYear = computed(() => new Date().getFullYear());
</script>
