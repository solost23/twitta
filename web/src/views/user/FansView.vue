<template>
  <div>
    <el-tabs v-model="tab">
      <el-tab-pane label="我的粉丝" name="fans">
        <div v-for="u in fans" :key="u.userId" class="user-row" @click="router.push(`/user/${u.userId}`)">
          <el-avatar :src="ossUrl(u.avatar)" size="small" />
          <span class="intro">{{ u.introduce || '暂无简介' }}</span>
        </div>
        <el-empty v-if="!fans.length" description="暂无粉丝" />
      </el-tab-pane>
      <el-tab-pane label="我的关注" name="following">
        <div v-for="u in following" :key="u.userId" class="user-row">
          <el-avatar :src="ossUrl(u.avatar)" size="small" @click="router.push(`/user/${u.userId}`)" style="cursor:pointer" />
          <span class="intro" style="flex:1">{{ u.introduce || '暂无简介' }}</span>
          <el-button size="small" @click="unfollow(u.userId)">取消关注</el-button>
        </div>
        <el-empty v-if="!following.length" description="还没有关注任何人" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { socialApi, type FanItem } from '@/api/social'
import { ElMessage } from 'element-plus'

const router = useRouter()
const tab = ref('fans')
const fans = ref<FanItem[]>([])
const following = ref<FanItem[]>([])

function ossUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}

async function unfollow(id: string) {
  await socialApi.unfollow(id)
  ElMessage.success('已取消关注')
  following.value = following.value.filter(u => u.userId !== id)
}

onMounted(async () => {
  const [f, w] = await Promise.all([socialApi.fanList(), socialApi.followingList()])
  fans.value = f
  following.value = w
})
</script>

<style scoped>
.user-row { display: flex; align-items: center; gap: 10px; padding: 10px 0; border-bottom: 1px solid #f0f0f0; cursor: pointer; }
.intro { font-size: 13px; color: #555; }
</style>
