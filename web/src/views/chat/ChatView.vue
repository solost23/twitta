<template>
  <div class="chat-page">
    <el-card class="chat-card">
      <template #header>
        <div class="chat-header">
          <el-button :icon="ArrowLeft" link @click="router.back()" />
          <el-avatar :src="targetUser?.avatar" size="small" />
          <span>{{ targetUser?.nickname || targetUser?.username || targetId }}</span>
          <el-tag :type="wsStatus === 'open' ? 'success' : 'danger'" size="small">
            {{ wsStatus === 'open' ? '已连接' : '未连接' }}
          </el-tag>
        </div>
      </template>

      <!-- 历史消息加载 -->
      <div class="load-more" v-if="hasMore">
        <el-button size="small" link @click="loadMore">加载更多历史消息</el-button>
      </div>

      <!-- 消息列表 -->
      <div class="messages" ref="msgContainer">
        <div
          v-for="(msg, i) in messages"
          :key="i"
          :class="['msg-row', isMine(msg) ? 'mine' : 'theirs']"
        >
          <div class="bubble">
            <div class="text">{{ msg.content || msg.msg }}</div>
            <div class="time">{{ msg.createdAt }}</div>
          </div>
        </div>
      </div>

      <!-- 输入框 -->
      <div class="input-bar">
        <el-input
          v-model="inputText"
          placeholder="输入消息..."
          @keyup.enter="send"
          :disabled="wsStatus !== 'open'"
        />
        <el-button type="primary" :disabled="wsStatus !== 'open'" @click="send">发送</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { chatApi } from '@/api/chat'
import { authApi, type UserDetail } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { ArrowLeft } from '@element-plus/icons-vue'

interface DisplayMsg {
  fromId?: string; userId?: string
  content?: string; msg?: string; createdAt: string
}

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const targetId = route.params.id as string
const myId = auth.user?.id || ''

const messages = ref<DisplayMsg[]>([])
const inputText = ref('')
const wsStatus = ref<'connecting' | 'open' | 'closed'>('connecting')
const msgContainer = ref<HTMLElement | null>(null)
const hasMore = ref(false)
const page = ref(1)
const targetUser = ref<UserDetail | null>(null)

let ws: WebSocket | null = null

// ObjectID("abc123") 和纯 hex "abc123" 都能正确比较
function normalizeId(id: string): string {
  const m = id.match(/^ObjectID\("([a-f0-9]+)"\)$/i)
  return m ? m[1] : id
}

function isMine(msg: DisplayMsg): boolean {
  const uid = normalizeId(msg.fromId || msg.userId || '')
  return uid === normalizeId(myId)
}

function buildWsUrl() {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  const host = location.host
  const token = localStorage.getItem('token') || ''
  return `${proto}://${host}/api/twitta/chats/${targetId}/ws?token=${token}`
}

function connectWs() {
  ws = new WebSocket(buildWsUrl())
  wsStatus.value = 'connecting'

  ws.onopen = () => { wsStatus.value = 'open' }
  ws.onclose = () => { wsStatus.value = 'closed' }
  ws.onerror = () => { wsStatus.value = 'closed' }

  ws.onmessage = e => {
    const msg = JSON.parse(e.data) as DisplayMsg
    messages.value.push(msg)
    scrollBottom()
  }
}

function send() {
  if (!inputText.value.trim() || !ws || ws.readyState !== WebSocket.OPEN) return
  ws.send(JSON.stringify({ content: inputText.value }))
  inputText.value = ''
}

async function loadHistory(p = 1) {
  const res = await chatApi.history(targetId, p, 20)
  const older = res.records.map(r => ({ userId: r.userId, msg: r.msg, createdAt: r.createdAt }))
  messages.value = [...older, ...messages.value]
  hasMore.value = p < res.pages
  page.value = p
}

async function loadMore() {
  await loadHistory(page.value + 1)
}

function scrollBottom() {
  nextTick(() => {
    if (msgContainer.value) msgContainer.value.scrollTop = msgContainer.value.scrollHeight
  })
}

onMounted(async () => {
  const [_, detail] = await Promise.allSettled([
    loadHistory(1),
    authApi.userDetail(targetId)
  ])
  if (detail.status === 'fulfilled') targetUser.value = detail.value
  scrollBottom()
  connectWs()
})

onUnmounted(() => { ws?.close() })
</script>

<style scoped>
.chat-page { height: calc(100vh - 96px); display: flex; flex-direction: column; }
.chat-card { flex: 1; display: flex; flex-direction: column; height: 100%; }
.chat-card :deep(.el-card__body) { flex: 1; display: flex; flex-direction: column; padding: 0; overflow: hidden; }
.chat-header { display: flex; align-items: center; gap: 8px; }
.load-more { text-align: center; padding: 8px; }
.messages {
  flex: 1; overflow-y: auto; padding: 16px;
  display: flex; flex-direction: column; gap: 8px;
}
.msg-row { display: flex; }
.msg-row.mine { justify-content: flex-end; }
.msg-row.theirs { justify-content: flex-start; }
.bubble { max-width: 60%; }
.mine .bubble .text {
  background: #1da1f2; color: #fff;
  padding: 8px 12px; border-radius: 16px 16px 4px 16px;
}
.theirs .bubble .text {
  background: #f0f0f0; color: #333;
  padding: 8px 12px; border-radius: 16px 16px 16px 4px;
}
.time { font-size: 11px; color: #bbb; margin-top: 2px; text-align: right; }
.input-bar { display: flex; gap: 8px; padding: 12px 16px; border-top: 1px solid #e4e7ed; }
.input-bar .el-input { flex: 1; }
</style>
