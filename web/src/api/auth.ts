import http from '@/utils/http'

export interface LoginForm { username: string; password: string; platform: string }
export interface RegisterForm {
  username: string; password: string; nickname?: string
  mobile?: string; email: string; avatar?: string; introduce?: string
}
export interface UserDetail {
  id: string; username: string; nickname: string; avatar: string
  introduce: string; whatCount: number; fansCount: number; createdAt: string
}
export interface LoginResponse {
  id: string; username: string; nickname: string; avatar: string
  role: string; token: string; isFirstLogin: number
}

export const authApi = {
  login: (data: LoginForm) => http.post<any, LoginResponse>('/login', data),
  register: (data: RegisterForm) => http.post<any, void>('/register', data),
  uploadAvatar: (file: File) => {
    const fd = new FormData(); fd.append('file', file)
    return http.post<any, string>('/register/avatar', fd)
  },
  logout: () => http.post<any, void>('/users/logout'),
  userDetail: (id: string) => http.get<any, UserDetail>(`/users/${id}`),
  userUpdate: (data: Partial<RegisterForm>) => http.put<any, void>('/users', data),
  userSearch: (keyword: string, page = 1, size = 10) =>
    http.get<any, { records: UserDetail[]; total: number }>('/users/search', { params: { keyword, page, size } })
}
