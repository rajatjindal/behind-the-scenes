import { ofetch } from 'ofetch'

const runtimeConfig = useRuntimeConfig().public.baseURL
const baseURL = runtimeConfig

export const myfetch = ofetch.create({
    baseURL: baseURL,
    retry: 0,
    async onRequest({ request, options }) {
        console.log('[fetch request]', `[${options.method} ${request}]`)
    },
    onRequestError: function ({ request, options, error }) { },
    onResponse: function ({ request, response, options }) { },
    onResponseError: function ({ request, response, options, error }) {
        console.error('[fetch response error]', response.status, response.body, error);
    },
})