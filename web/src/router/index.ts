import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/feeds'
  },
  {
    path: '/torrents/unimported',
    name: 'unimported-torrents',
    component: () => import('@/views/UnimportedTorrentsView.vue')
  },
  {
    path: '/torrents/manual',
    name: 'manual-torrents',
    component: () => import('@/views/ManualTorrentsView.vue')
  },
  {
    path: '/feeds',
    name: 'feeds',
    component: () => import('@/views/FeedsView.vue')
  },
  {
    path: '/libraries',
    name: 'libraries',
    component: () => import('@/views/LibrariesView.vue')
  },
  {
    path: '/feed-filters',
    name: 'feed-filters',
    component: () => import('@/views/FeedFiltersView.vue')
  },
  {
    path: '/feed-author-filters',
    name: 'feed-author-filters',
    component: () => import('@/views/FeedAuthorFiltersView.vue')
  },
  {
    path: '/notifiers',
    name: 'notifiers',
    component: () => import('@/views/NotifiersView.vue')
  },
  {
    path: '/filter-keys',
    name: 'filter-keys',
    component: () => import('@/views/FilterKeysView.vue')
  },
  {
    path: '/filter-operators',
    name: 'filter-operators',
    component: () => import('@/views/FilterOperatorsView.vue')
  },
  {
    path: '/notification-types',
    name: 'notification-types',
    component: () => import('@/views/NotificationTypesView.vue')
  },
  {
    path: '/feed-filter-set-types',
    name: 'feed-filter-set-types',
    component: () => import('@/views/FeedFilterSetTypesView.vue')
  },
  {
    path: '/torrent-categories',
    name: 'torrent-categories',
    component: () => import('@/views/TorrentCategoriesView.vue')
  },
  {
    path: '/book-types',
    name: 'book-types',
    component: () => import('@/views/BookTypesView.vue')
  },
  {
    path: '/authors',
    name: 'authors',
    component: () => import('@/views/AuthorsView.vue')
  },
  {
    path: '/subscription-items',
    name: 'subscription-items',
    component: () => import('@/views/SubscriptionItemsView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
