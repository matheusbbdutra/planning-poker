import type { UseFetchOptions } from '#app';

/**
 * Wrapper para o useFetch do Nuxt que configura automaticamente a baseURL da API.
 * Isso centraliza a lógica de comunicação com o backend.
 *
 * @param path O caminho do endpoint da API (ex: '/rooms').
 * @param options Opções adicionais para o useFetch.
 */
export const useApi = <T>(path: string, options: UseFetchOptions<T> = {}) => {
  const config = useRuntimeConfig();

  const defaults: UseFetchOptions<T> = {
    baseURL: config.public.apiBase,
  };

  return useFetch(path, { ...defaults, ...options });
};

