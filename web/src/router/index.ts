import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: () => import('@/views/auth/LoginView.vue') },
    { path: '/register', component: () => import('@/views/auth/RegisterView.vue') },
    {
      path: '/',
      component: () => import('@/views/LayoutView.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/home' },
        { path: 'home', component: () => import('@/views/tweet/HomeView.vue') },
        { path: 'search', component: () => import('@/views/tweet/SearchView.vue') },
        { path: 'favorites', component: () => import('@/views/tweet/FavoritesView.vue') },
        { path: 'own', component: () => import('@/views/tweet/OwnTweetsView.vue') },
        { path: 'profile', component: () => import('@/views/user/ProfileView.vue') },
        { path: 'user/:id', component: () => import('@/views/user/UserDetailView.vue') },
        { path: 'friends', component: () => import('@/views/friend/FriendsView.vue') },
        { path: 'fans', component: () => import('@/views/user/FansView.vue') },
        { path: 'chat/:id', component: () => import('@/views/chat/ChatView.vue') }
      ]
    },
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ]
})

router.beforeEach(to => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isLoggedIn) return '/login'
})

export default router
