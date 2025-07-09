export function genAuthorization(value: string | undefined) {
  if (!value) {
    return undefined
  }

  return {
    Authorization: `Bearer ${value}`,
  }
}
