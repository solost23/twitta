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
      <el-button link :class="{ liked }" @click="toggleThumb">
        <el-icon><Pointer /></el-icon> {{ localThumbCount }}
      </el-button>
      <el-button link @click="toggleComments">
        <el-icon><ChatLineRound /></el-icon> {{ localCommentCount }}
      </el-button>
      <el-button link @click="toggleFavorite">
        <el-icon><Star /></el-icon> {{ favorited ? '已收藏' : '收藏' }}
      </el-button>
    </div>

    <!-- 评论区 -->
    <div v-if="showComments" class="comments">
      <div class="comment-input">
        <el-input
          v-model="commentText"
          placeholder="写评论..."
          size="small"
          style="flex:1"
          @keyup.enter="submitComment"
        />
        <el-button size="small" type="primary" :loading="submitting" @click="submitComment">发送</el-button>
      </div>
      <div v-if="loading" class="loading-tip">加载中...</div>
      <el-empty v-else-if="!comments.length" description="暂无评论" :image-size="40" />
      <CommentItem
        v-for="c in comments"
        :key="c.id"
        :comment="c"
        :tweet-id="tweet.id"
        @replied="loadComments"
        @deleted="onCommentDeleted"
      />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
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
const loading = ref(false)
const submitting = ref(false)
const liked = ref(false)
const favorited = ref(false)
const localThumbCount = ref(props.tweet.thumbCount)
const localCommentCount = ref(props.tweet.commentCount)

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

async function toggleThumb() {
  if (liked.value) {
    await tweetApi.unthumb(props.tweet.id)
    liked.value = false
    localThumbCount.value--
  } else {
    await tweetApi.thumb(props.tweet.id)
    liked.value = true
    localThumbCount.value++
  }
}

async function toggleFavorite() {
  if (favorited.value) {
    await tweetApi.unfavorite(props.tweet.id)
    favorited.value = false
    ElMessage.success('已取消收藏')
  } else {
    await tweetApi.favorite(props.tweet.id)
    favorited.value = true
    ElMessage.success('已收藏')
  }
}

async function toggleComments() {
  showComments.value = !showComments.value
  if (showComments.value && !comments.value.length) await loadComments()
}

async function loadComments() {
  loading.value = true
  try {
    const res = await tweetApi.commentList(props.tweet.id)
    comments.value = res.records ?? []
  } finally {
    loading.value = false
  }
}

async function submitComment() {
  if (!commentText.value.trim()) return
  submitting.value = true
  try {
    await tweetApi.comment(props.tweet.id, commentText.value)
    commentText.value = ''
    localCommentCount.value++
    await loadComments()
  } finally {
    submitting.value = false
  }
}

function onCommentDeleted() {
  localCommentCount.value = Math.max(0, localCommentCount.value - 1)
  loadComments()
}
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
.actions .liked { color: #1da1f2; }
.comments { margin-top: 12px; border-top: 1px solid #f0f0f0; padding-top: 8px; }
.comment-input { display: flex; gap: 8px; margin-bottom: 12px; }
.loading-tip { font-size: 13px; color: #bbb; text-align: center; padding: 8px; }
</style>
