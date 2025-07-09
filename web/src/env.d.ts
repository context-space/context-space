interface ImportMetaEnv {
  readonly NEXT_PUBLIC_SUPABASE_URL: string
  readonly NEXT_PUBLIC_SUPABASE_ANON_KEY: string
  readonly NEXT_PUBLIC_BASE_API_URL: string
  readonly NEXT_PUBLIC_BASE_URL: string
  readonly OPENAI_API_KEY: string
  readonly OPENAI_BASE_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
