<template>
  <div>
    <!-- 发推框 -->
    <el-card class="compose-card">
      <el-input v-model="newTweet.title" placeholder="标题（可选）" style="margin-bottom:8px" />
      <el-input v-model="newTweet.content" type="textarea" :rows="3" placeholder="有什么新鲜事？" />
      <div class="compose-actions">
        <el-upload action="#" :before-upload="uploadMedia" :show-file-list="false" accept="image/*,video/*" multiple>
          <el-button size="small" :icon="Picture">添加图片</el-button>
        </el-upload>
        <div class="image-preview">
          <el-image v-for="(url, i) in newTweet.images" :key="i" :src="url" style="width:60px;height:60px;object-fit:cover;border-radius:4px" />
        </div>
        <el-button type="primary" size="small" :loading="sending" @click="sendTweet">发布</el-button>
      </div>
    </el-card>

    <!-- 推文列表 -->
    <TweetCard
      v-for="tweet in tweets"
      :key="tweet.id"
      :tweet="tweet"
      @deleted="loadTweets"
      @thumb="handleThumb"
    />

    <el-pagination
      v-if="total > pageSize"
      layout="prev, pager, next"
      :total="total"
      :page-size="pageSize"
      :current-page="page"
      @current-change="p => { page = p; loadTweets() }"
      style="margin-top:16px;justify-content:center;display:flex"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { tweetApi, type Tweet } from '@/api/tweet'
import { ElMessage } from 'element-plus'
import { Picture } from '@element-plus/icons-vue'
import TweetCard from '@/components/TweetCard.vue'

const tweets = ref<Tweet[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const sending = ref(false)
const newTweet = reactive({ title: '', content: '', images: [] as string[] })

async function loadTweets() {
  const res = await tweetApi.list(page.value, pageSize)
  tweets.value = res.records
  total.value = res.total
}

async function uploadMedia(file: File) {
  const url = await tweetApi.uploadStatic(file)
  newTweet.images.push(url)
  return false
}

async function sendTweet() {
  if (!newTweet.content.trim()) { ElMessage.warning('内容不能为空'); return }
  sending.value = true
  try {
    await tweetApi.send({ title: newTweet.title, content: newTweet.content, images: newTweet.images })
    newTweet.title = ''; newTweet.content = ''; newTweet.images = []
    ElMessage.success('发布成功')
    loadTweets()
  } finally {
    sending.value = false
  }
}

async function handleThumb(tweetId: string, liked: boolean) {
  if (liked) await tweetApi.thumb(tweetId)
  else await tweetApi.unthumb(tweetId)
  loadTweets()
}

onMounted(loadTweets)
</script>

<style scoped>
.compose-card { margin-bottom: 16px; }
.compose-actions { display: flex; align-items: center; gap: 8px; margin-top: 8px; }
.image-preview { display: flex; gap: 4px; flex: 1; flex-wrap: wrap; }
</style>
