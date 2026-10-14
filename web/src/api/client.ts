// 统一 HTTP 客户端：注入 JWT、解包统一响应体 {code,message,data,traceId}、
// 处理 401 跳转登录、把非 0 code 转成 ApiError 抛给调用方。

export const TOKEN_KEY = 'ra_token'

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}
export function setToken(t: string): void {
  localStorage.setItem(TOKEN_KEY, t)
}
export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

export class ApiError extends Error {
  code: number
  constructor(code: number, message: string) {
    super(message)
    this.code = code
  }
}

type UnauthorizedHandler = () => void
let unauthorizedHandler: UnauthorizedHandler | null = null

// 注册 401 时的全局回调（用于跳转登录页）。
export function onUnauthorized(fn: UnauthorizedHandler): void {
  unauthorizedHandler = fn
}

interface Envelope<T> {
  code: number
  message: string
  data: T
  traceId?: string
}

export type Query = Record<string, string | number | boolean | undefined | null>

async function request<T>(
  method: string,
  path: string,
  opts: { query?: Query; body?: unknown } = {},
): Promise<T> {
  let url = '/api/v1' + path
  if (opts.query) {
    const qs = new URLSearchParams()
    for (const [k, v] of Object.entries(opts.query)) {
      if (v !== undefined && v !== null && v !== '') qs.append(k, String(v))
    }
    const s = qs.toString()
    if (s) url += '?' + s
  }

  const headers: Record<string, string> = {}
  const token = getToken()
  if (token) headers.Authorization = 'Bearer ' + token
  if (opts.body !== undefined) headers['Content-Type'] = 'application/json'

  let res: Response
  try {
    res = await fetch(url, {
      method,
      headers,
      body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
    })
  } catch {
    throw new ApiError(-1, 'network error')
  }

  let body: Envelope<T>
  try {
    body = (await res.json()) as Envelope<T>
  } catch {
    throw new ApiError(res.status, 'invalid response')
  }

  if (res.status === 401 || body.code === 10002) {
    clearToken()
    unauthorizedHandler?.()
    throw new ApiError(10002, body.message || 'unauthorized')
  }
  if (body.code !== 0) {
    throw new ApiError(body.code, body.message || 'request failed')
  }
  return body.data
}

export const http = {
  get: <T>(path: string, query?: Query) => request<T>('GET', path, { query }),
  post: <T>(path: string, body?: unknown) => request<T>('POST', path, { body }),
  put: <T>(path: string, body?: unknown) => request<T>('PUT', path, { body }),
  del: <T>(path: string) => request<T>('DELETE', path),
}
