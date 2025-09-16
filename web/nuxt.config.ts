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
    nitro: {
        storage: {
            redis: {
                driver: 'redis',
                host: '127.0.0.1',
                port: 6379,
            }
        },
    },
    runtimeConfig: {
        redis: {
            host: '127.0.0.1',
            port: 6379
        }
    }
})
