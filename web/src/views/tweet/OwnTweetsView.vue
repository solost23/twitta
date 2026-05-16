<template>
  <div>
    <h3 style="margin-bottom:16px">我的推文</h3>
    <TweetCard v-for="t in tweets" :key="t.id" :tweet="t" @deleted="load" @thumb="() => {}" />
    <el-empty v-if="!tweets.length" description="还没有发过推文" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { tweetApi, type Tweet } from '@/api/tweet'
import TweetCard from '@/components/TweetCard.vue'

const tweets = ref<Tweet[]>([])
async function load() {
  const res = await tweetApi.ownList()
  tweets.value = res.records
}
onMounted(load)
</script>
