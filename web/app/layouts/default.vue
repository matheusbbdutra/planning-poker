<template>
  <div class="flex min-h-screen flex-col">
    <NavbarSection />
    <main class="flex-1">
      <slot />
    </main>
    <FooterSection
      authorName="Matheus Dutra"
      projectLink="https://github.com/matheusbbdutra/planning-poker"
      contactEmail="contato@mdutra.dev"
    />
  </div>
</template>
<script setup lang="ts">
import NavbarSection from "@/components/NavbarSection.vue";
import FooterSection from "@/components/FooterSection.vue";
</script>
