<template>
  <div>
    <h3 style="margin-bottom:16px">我的收藏</h3>
    <TweetCard v-for="t in tweets" :key="t.id" :tweet="t" @deleted="load" />
    <el-empty v-if="!tweets.length" description="暂无收藏" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { tweetApi, type Tweet } from '@/api/tweet'
import TweetCard from '@/components/TweetCard.vue'

const tweets = ref<Tweet[]>([])
async function load() {
  const res = await tweetApi.favoriteList()
  tweets.value = res.records
}
onMounted(load)
</script>
