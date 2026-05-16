import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useNotifyStore = defineStore('notify', () => {
  const unread = ref(0)

  function increment() { unread.value++ }
  function clear() { unread.value = 0 }

  return { unread, increment, clear }
})
