import http from '@/utils/http'

export interface Tweet {
  id: string; userId: string; username: string; avatar: string
  title: string; content: string; images: string[]
  createdAt: string; thumbCount: number; commentCount: number
}
export interface TweetList { records: Tweet[]; total: number; pages: number; size: number; current: number }
export interface Comment {
  id: string; pid: string; userId: string; username: string
  avatar: string; content: string; children: Comment[]
}

export const tweetApi = {
  list: (page = 1, size = 10) =>
    http.get<any, TweetList>('/tweets', { params: { page, size } }),
  ownList: () => http.get<any, TweetList>('/tweets/own'),
  send: (data: { title: string; content: string; images?: string[] }) =>
    http.post<any, void>('/tweets', data),
  uploadStatic: (file: File) => {
    const fd = new FormData(); fd.append('file', file)
    return http.post<any, string>('/tweets/static', fd)
  },
  delete: (id: string) => http.delete<any, void>(`/tweets/${id}`),
  search: (keyword: string, page = 1, size = 10) =>
    http.get<any, TweetList>('/tweets/search', { params: { keyword, page, size } }),
  favorite: (id: string) => http.post<any, void>('/tweets/favorite', { id }),
  unfavorite: (id: string) => http.delete<any, void>(`/tweets/favorite/${id}`),
  favoriteList: () => http.get<any, TweetList>('/tweets/favorite'),
  commentList: (tweetId: string, page = 1, size = 20) =>
    http.get<any, { records: Comment[] }>(`/comments/${tweetId}`, { params: { page, size } }),
  thumb: (tweetId: string) => http.post<any, void>(`/comments/${tweetId}/thumb`),
  unthumb: (tweetId: string) => http.delete<any, void>(`/comments/${tweetId}/thumb`),
  comment: (tweetId: string, content: string, parentId = '') =>
    http.post<any, void>(`/comments/${tweetId}`, { content, parentId }),
  deleteComment: (commentId: string) => http.delete<any, void>(`/comments/${commentId}`)
}
