import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: { enabled: true },
    modules: ['@pinia/nuxt'],
    css: ['~/assets/css/main.css'],
    vite: {
      plugins: [
          tailwindcss(),
      ],
    },
    runtimeConfig: {
      public: {
        apiBase: process.env.API_URL || 'http://localhost:8080/v1'
      }
    }
    
})
