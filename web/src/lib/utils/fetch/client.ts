import { $fetch } from "ofetch"
import { baseAPIURL } from "@/config"
import { getClientToken } from "@/lib/supabase/client"
import { genAuthorization } from "./shared"

export const clientRemoteFetch = $fetch.create({
  baseURL: baseAPIURL,
  timeout: 5000,
  headers: genAuthorization(await getClientToken()),
  // retryStatusCodes: [],
  // retry: 3,
})
