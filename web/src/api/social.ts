import http from '@/utils/http'

export interface FriendApplication {
  userId: string; username: string; avatar: string
  content: string; type: number; createdAt: string
}
export interface FriendItem {
  userId: string; username: string; nickname: string; avatar: string; introduce: string
}
export interface FanItem { userId: string; avatar: string; introduce: string }

export const socialApi = {
  fanList: () => http.get<any, FanItem[]>('/fans'),
  followingList: () => http.get<any, FanItem[]>('/fans/what'),
  follow: (id: string) => http.post<any, void>(`/fans/${id}`),
  unfollow: (id: string) => http.delete<any, void>(`/fans/${id}`),
  friendList: () => http.get<any, FriendItem[]>('/friends/list'),
  applicationList: () => http.get<any, FriendApplication[]>('/friends'),
  sendApplication: (userId: string, content: string) =>
    http.post<any, void>('/friends', { userId, content }),
  acceptApplication: (id: string) => http.put<any, void>(`/friends/${id}/accept`),
  rejectApplication: (id: string) => http.put<any, void>(`/friends/${id}/reject`),
  deleteFriend: (id: string) => http.delete<any, void>(`/friends/${id}`)
}
