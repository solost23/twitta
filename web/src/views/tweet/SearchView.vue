<template>
  <div>
    <el-input v-model="keyword" placeholder="搜索用户或推文..." clearable @keyup.enter="search" style="margin-bottom:16px">
      <template #append><el-button :icon="Search" @click="search" /></template>
    </el-input>

    <el-tabs v-if="searched" v-model="tab">
      <el-tab-pane label="推文" name="tweet">
        <TweetCard v-for="t in tweetResults" :key="t.id" :tweet="t" @deleted="search" />
        <el-empty v-if="!tweetResults.length" description="无结果" />
      </el-tab-pane>
      <el-tab-pane label="用户" name="user">
        <el-card v-for="u in userResults" :key="u.id" class="user-card" @click="router.push(`/user/${u.id}`)">
          <el-avatar :src="ossUrl(u.avatar)" />
          <div class="info">
            <div class="name">{{ u.nickname || u.username }}</div>
            <div class="intro">{{ u.introduce }}</div>
          </div>
        </el-card>
        <el-empty v-if="!userResults.length" description="无结果" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { tweetApi, type Tweet } from '@/api/tweet'
import { authApi, type UserDetail } from '@/api/auth'
import { Search } from '@element-plus/icons-vue'
import TweetCard from '@/components/TweetCard.vue'

const router = useRouter()
const keyword = ref('')
const tab = ref('tweet')
const searched = ref(false)
const tweetResults = ref<Tweet[]>([])
const userResults = ref<UserDetail[]>([])

function ossUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}

async function search() {
  if (!keyword.value.trim()) return
  searched.value = true
  const [tr, ur] = await Promise.all([
    tweetApi.search(keyword.value),
    authApi.userSearch(keyword.value)
  ])
  tweetResults.value = tr.records
  userResults.value = ur.records
}
</script>

<style scoped>
.user-card { display: flex; align-items: center; gap: 12px; margin-bottom: 8px; cursor: pointer; }
.user-card:hover { background: #f5f7fa; }
.info .name { font-weight: 600; }
.info .intro { font-size: 12px; color: #999; }
</style>
