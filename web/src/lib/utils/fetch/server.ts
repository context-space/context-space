import type { BaseRemoteAPIResponse } from "@/typings"
import { $fetch } from "ofetch"
import { baseAPIURL } from "@/config"

export const serverRemoteFetch = $fetch.create({
  baseURL: baseAPIURL,
  // timeout: 5000,
  // retryStatusCodes: [],
  // retry: 3,
})

export async function remoteFetchWithErrorHandling<T>(
  url: string,
  options: RequestInit = {},
): Promise<T> {
  const response = (await serverRemoteFetch(url, options)) as BaseRemoteAPIResponse<T>
  if (!response.success) {
    throw new Error(response.message)
  }
  return response.data
}
