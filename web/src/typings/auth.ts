export interface AuthActionResult {
  success: boolean
  error?: string
  redirectTo?: string
}

export interface LoginResult extends AuthActionResult {
  redirectTo?: string
}

export interface SendCodeResult {
  success: boolean
  error?: string
  message?: string
}
