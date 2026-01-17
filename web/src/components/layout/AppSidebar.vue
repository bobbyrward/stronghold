<script setup lang="ts">
import { useUiStore } from '@/stores/ui'

const uiStore = useUiStore()
</script>

<template>
    <aside class="app-sidebar" :class="{ collapsed: uiStore.sidebarCollapsed }">
        <nav class="sidebar-nav">
            <div class="section-title">Main</div>

            <router-link to="/feeds" class="nav-link">
                <i class="bi bi-rss"></i>
                <span class="nav-text">Feeds</span>
            </router-link>

            <router-link to="/feed-filters" class="nav-link">
                <i class="bi bi-funnel"></i>
                <span class="nav-text">Feed Filters</span>
            </router-link>

            <router-link to="/feed-author-filters" class="nav-link">
                <i class="bi bi-person-badge"></i>
                <span class="nav-text">Author Filters</span>
            </router-link>

            <router-link to="/notifiers" class="nav-link">
                <i class="bi bi-bell"></i>
                <span class="nav-text">Notifiers</span>
            </router-link>

            <div class="section-title">Subscriptions</div>

            <router-link to="/authors" class="nav-link">
                <i class="bi bi-person-lines-fill"></i>
                <span class="nav-text">Authors</span>
            </router-link>

            <router-link to="/subscription-items" class="nav-link">
                <i class="bi bi-download"></i>
                <span class="nav-text">Download History</span>
            </router-link>

            <div class="section-title">Torrents</div>

            <router-link to="/torrents/unimported" class="nav-link">
                <i class="bi bi-cloud-arrow-down"></i>
                <span class="nav-text">Unimported</span>
            </router-link>

            <router-link to="/torrents/manual" class="nav-link">
                <i class="bi bi-cloud-arrow-down"></i>
                <span class="nav-text">Manual Intervention</span>
            </router-link>


            <div class="section-title">Reference Data</div>

            <router-link to="/filter-keys" class="nav-link">
                <i class="bi bi-key"></i>
                <span class="nav-text">Filter Keys</span>
            </router-link>

            <router-link to="/filter-operators" class="nav-link">
                <i class="bi bi-sliders"></i>
                <span class="nav-text">Filter Operators</span>
            </router-link>

            <router-link to="/notification-types" class="nav-link">
                <i class="bi bi-envelope"></i>
                <span class="nav-text">Notification Types</span>
            </router-link>

            <router-link to="/feed-filter-set-types" class="nav-link">
                <i class="bi bi-collection"></i>
                <span class="nav-text">Filter Set Types</span>
            </router-link>

            <router-link to="/torrent-categories" class="nav-link">
                <i class="bi bi-folder"></i>
                <span class="nav-text">Torrent Categories</span>
            </router-link>
        </nav>
    </aside>
</template>
