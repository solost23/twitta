import http from '@/utils/http'

export interface ChatMessage {
  userId: string; msg: string; createdAt: string
}
export interface ChatList {
  records: ChatMessage[]; total: number; pages: number; size: number; current: number
}

export const chatApi = {
  history: (targetId: string, page = 1, size = 20) =>
    http.get<any, ChatList>(`/chats/${targetId}`, { params: { page, size } })
}
