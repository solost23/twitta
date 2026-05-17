<template>
  <div class="comment-item">
    <el-avatar :src="ossUrl(comment.avatar)" size="small" class="avatar" />
    <div class="body">
      <div class="meta-row">
        <span class="name">{{ comment.username }}</span>
        <span class="time" v-if="comment.createdAt">{{ formatTime(comment.createdAt) }}</span>
        <el-button v-if="isOwn" link size="small" class="del-btn" @click="deleteComment">删除</el-button>
        <el-button link size="small" class="reply-btn" @click="showReply = !showReply">回复</el-button>
      </div>
      <div class="text">{{ comment.content }}</div>

      <!-- 回复输入框 -->
      <div v-if="showReply" class="reply-input">
        <el-input v-model="replyText" size="small" placeholder="回复..." style="flex:1" @keyup.enter="submitReply" />
        <el-button size="small" type="primary" @click="submitReply">发送</el-button>
        <el-button size="small" @click="showReply = false">取消</el-button>
      </div>

      <!-- 子评论 -->
      <div v-if="comment.children?.length" class="children">
        <CommentItem
          v-for="c in comment.children"
          :key="c.id"
          :comment="c"
          :tweet-id="tweetId"
          @replied="emit('replied')"
          @deleted="emit('deleted')"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { tweetApi, type Comment } from '@/api/tweet'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const props = defineProps<{ comment: Comment; tweetId: string }>()
const emit = defineEmits<{ replied: []; deleted: [] }>()

const auth = useAuthStore()
const showReply = ref(false)
const replyText = ref('')

const isOwn = computed(() => {
  const myId = auth.user?.id || ''
  const cId = props.comment.userId || ''
  return myId === cId || myId === cId.replace(/^ObjectID\("([a-f0-9]+)"\)$/i, '$1')
})

function ossUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}

function formatTime(t: string) {
  if (!t) return ''
  const d = new Date(t)
  if (isNaN(d.getTime())) return t
  const now = new Date()
  return d.toDateString() === now.toDateString()
    ? d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    : d.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

async function submitReply() {
  if (!replyText.value.trim()) return
  await tweetApi.comment(props.tweetId, replyText.value, props.comment.id)
  replyText.value = ''
  showReply.value = false
  emit('replied')
}

async function deleteComment() {
  await ElMessageBox.confirm('确认删除这条评论？', '提示', { type: 'warning' })
  await tweetApi.deleteComment(props.comment.id)
  ElMessage.success('已删除')
  emit('deleted')
}
</script>

<style scoped>
.comment-item { display: flex; gap: 8px; margin-bottom: 10px; }
.avatar { flex-shrink: 0; }
.body { flex: 1; min-width: 0; }
.meta-row { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; }
.name { font-weight: 600; font-size: 13px; }
.time { font-size: 11px; color: #bbb; }
.del-btn, .reply-btn { font-size: 11px; padding: 0; height: auto; color: #bbb; }
.del-btn:hover { color: #f56c6c; }
.reply-btn:hover { color: #1da1f2; }
.text { font-size: 13px; color: #333; margin-top: 2px; word-break: break-word; }
.reply-input { display: flex; gap: 6px; margin-top: 6px; }
.children { margin-top: 8px; padding-left: 12px; border-left: 2px solid #e4e7ed; }
</style>
