<template>
  <el-card class="tweet-card">
    <div class="header">
      <el-avatar :src="ossUrl(tweet.avatar)" size="small" @click="goUser" style="cursor:pointer" />
      <div class="meta">
        <span class="username" @click="goUser">{{ tweet.username }}</span>
        <span class="time">{{ tweet.createdAt }}</span>
      </div>
      <el-button v-if="isOwn" link :icon="Delete" @click="deleteTweet" />
    </div>
    <div v-if="tweet.title" class="title">{{ tweet.title }}</div>
    <div class="content">{{ tweet.content }}</div>
    <div v-if="tweet.images?.length" class="images">
      <el-image
        v-for="(img, i) in tweet.images" :key="i"
        :src="ossUrl(img)"
        :preview-src-list="tweet.images.map(ossUrl)"
        :initial-index="i"
        fit="cover"
        class="thumb-img"
      />
    </div>
    <div class="actions">
      <el-button link @click="emit('thumb', tweet.id, true)">
        <el-icon><Pointer /></el-icon> {{ tweet.thumbCount }}
      </el-button>
      <el-button link @click="showComments = !showComments">
        <el-icon><ChatLineRound /></el-icon> {{ tweet.commentCount }}
      </el-button>
      <el-button link @click="favorite">
        <el-icon><Star /></el-icon> 收藏
      </el-button>
    </div>

    <!-- 评论区 -->
    <div v-if="showComments" class="comments">
      <div class="comment-input">
        <el-input v-model="commentText" placeholder="写评论..." size="small" style="flex:1" />
        <el-button size="small" type="primary" @click="submitComment">发送</el-button>
      </div>
      <CommentItem v-for="c in comments" :key="c.id" :comment="c" />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { tweetApi, type Tweet, type Comment } from '@/api/tweet'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Pointer, ChatLineRound, Star } from '@element-plus/icons-vue'
import CommentItem from './CommentItem.vue'

const props = defineProps<{ tweet: Tweet }>()
const emit = defineEmits<{ deleted: []; thumb: [id: string, liked: boolean] }>()

const router = useRouter()
const auth = useAuthStore()
const showComments = ref(false)
const comments = ref<Comment[]>([])
const commentText = ref('')
const isOwn = computed(() => auth.user?.id === props.tweet.userId)

function ossUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}

function goUser() { router.push(`/user/${props.tweet.userId}`) }

async function deleteTweet() {
  await ElMessageBox.confirm('确认删除这条推文？', '提示', { type: 'warning' })
  await tweetApi.delete(props.tweet.id)
  ElMessage.success('已删除')
  emit('deleted')
}

async function favorite() {
  await tweetApi.favorite(props.tweet.id)
  ElMessage.success('已收藏')
}

async function loadComments() {
  const res = await tweetApi.commentList(props.tweet.id)
  comments.value = res.records
}

async function submitComment() {
  if (!commentText.value.trim()) return
  await tweetApi.comment(props.tweet.id, commentText.value)
  commentText.value = ''
  loadComments()
}

watch(showComments, v => { if (v) loadComments() })
</script>

<style scoped>
.tweet-card { margin-bottom: 12px; }
.header { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.meta { flex: 1; display: flex; flex-direction: column; }
.username { font-weight: 600; font-size: 14px; cursor: pointer; }
.username:hover { color: #1da1f2; }
.time { font-size: 12px; color: #999; }
.title { font-weight: 700; font-size: 15px; margin-bottom: 4px; }
.content { font-size: 14px; line-height: 1.6; white-space: pre-wrap; }
.images { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 8px; }
.thumb-img { width: 120px; height: 120px; border-radius: 6px; cursor: pointer; }
.actions { display: flex; gap: 16px; margin-top: 8px; border-top: 1px solid #f0f0f0; padding-top: 8px; }
.comments { margin-top: 12px; border-top: 1px solid #f0f0f0; padding-top: 8px; }
.comment-input { display: flex; gap: 8px; margin-bottom: 8px; }
</style>
