import { $fetch } from "ofetch"

export const clientRemoteFetch = $fetch.create({
  baseURL: "/api/proxy",
})
