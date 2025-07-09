import { $fetch } from "ofetch"

export * from "./shared"

export const localFetch = $fetch.create({
  timeout: 2000,
  // retry: 3,
})
