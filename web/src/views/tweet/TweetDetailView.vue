<template>
  <div v-if="tweet">
    <el-page-header @back="router.back()" style="margin-bottom:16px" />
    <TweetCard :tweet="tweet" @deleted="router.back()" />
  </div>
  <div v-else-if="loading" style="text-align:center;padding:40px">
    <el-icon class="is-loading" :size="32"><Loading /></el-icon>
  </div>
  <el-empty v-else description="推文不存在" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { tweetApi, type Tweet } from '@/api/tweet'
import { Loading } from '@element-plus/icons-vue'
import TweetCard from '@/components/TweetCard.vue'

const route = useRoute()
const router = useRouter()
const tweet = ref<Tweet | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    tweet.value = await tweetApi.detail(route.params.id as string)
  } finally {
    loading.value = false
  }
})
</script>
