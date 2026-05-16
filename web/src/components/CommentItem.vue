<template>
  <div class="comment-item">
    <el-avatar :src="ossUrl(comment.avatar)" size="small" />
    <div class="body">
      <span class="name">{{ comment.username }}</span>
      <span class="text">{{ comment.content }}</span>
      <div v-if="comment.children?.length" class="children">
        <CommentItem v-for="c in comment.children" :key="c.id" :comment="c" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Comment } from '@/api/tweet'
defineProps<{ comment: Comment }>()
function ossUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}
</script>

<style scoped>
.comment-item { display: flex; gap: 8px; margin-bottom: 8px; }
.body { flex: 1; }
.name { font-weight: 600; font-size: 13px; margin-right: 6px; }
.text { font-size: 13px; color: #333; }
.children { margin-top: 6px; padding-left: 12px; border-left: 2px solid #e4e7ed; }
</style>
