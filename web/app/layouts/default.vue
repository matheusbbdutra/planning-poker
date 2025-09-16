<template>
  <div>
  <NavbarSection/>
  <main class="h-dvh">
    <slot/>
  </main>
    <FooterSection authorName="Matheus Dutra"/>
  </div>
</template>
<script setup lang="ts">
import NavbarSection from "@/components/NavbarSection.vue";
import FooterSection from "@/components/FooterSection.vue";
</script>