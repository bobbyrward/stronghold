<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import EventDetailDrawer from '@/components/activity/EventDetailDrawer.vue'
import type { EventLog, PaginatedEventLogResponse } from '@/types/api'

const toast = useToastStore()

const loading = ref(false)
const response = ref<PaginatedEventLogResponse | null>(null)
const selectedEvent = ref<EventLog | null>(null)
const autoRefresh = ref(false)
let autoRefreshTimer: ReturnType<typeof setInterval> | null = null

// Filter state
const filterCategory = ref('')
const filterSource = ref('')
const filterEventType = ref('')
const filterEntityType = ref('')
const filterQuery = ref('')
const filterFrom = ref('')
const filterTo = ref('')
const page = ref(1)
const perPage = ref(50)

// Debounce timer for text search
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const categoryColors: Record<string, string> = {
    download: 'bg-primary',
    import: 'bg-success',
    notification: 'bg-info',
    subscription: 'bg-warning text-dark',
    search: 'bg-secondary',
    feed: 'bg-dark',
    mutation: 'bg-danger',
}

const sourceColors: Record<string, string> = {
    'feedwatcher2': 'bg-primary',
    'discord-bot': 'bg-info',
    'api': 'bg-secondary',
    'ebook-importer': 'bg-success',
    'audiobook-importer': 'bg-warning text-dark',
    'author-subscription-importer': 'bg-danger',
}

function categoryBadgeClass(category: string): string {
    return categoryColors[category] || 'bg-secondary'
}

function sourceBadgeClass(source: string): string {
    return sourceColors[source] || 'bg-secondary'
}

function formatRelativeTime(iso: string): string {
    const now = Date.now()
    const then = new Date(iso).getTime()
    const diff = now - then

    const seconds = Math.floor(diff / 1000)
    if (seconds < 60) return `${seconds}s ago`
    const minutes = Math.floor(seconds / 60)
    if (minutes < 60) return `${minutes}m ago`
    const hours = Math.floor(minutes / 60)
    if (hours < 24) return `${hours}h ago`
    const days = Math.floor(hours / 24)
    return `${days}d ago`
}

function toLocalDatetimeValue(iso: string): string {
    if (!iso) return ''
    const d = new Date(iso)
    // Format as YYYY-MM-DDTHH:mm for datetime-local input
    const pad = (n: number) => n.toString().padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function fromLocalDatetimeValue(val: string): string {
    if (!val) return ''
    return new Date(val).toISOString()
}

async function loadData() {
    loading.value = true
    try {
        const params: Record<string, string> = {
            page: page.value.toString(),
            per_page: perPage.value.toString(),
        }
        if (filterCategory.value) params.category = filterCategory.value
        if (filterSource.value) params.source = filterSource.value
        if (filterEventType.value) params.event_type = filterEventType.value
        if (filterEntityType.value) params.entity_type = filterEntityType.value
        if (filterQuery.value) params.q = filterQuery.value
        if (filterFrom.value) params.from = filterFrom.value
        if (filterTo.value) params.to = filterTo.value

        response.value = await api.eventLogs.list(params)
    } catch {
        toast.error('Failed to load activity logs')
    } finally {
        loading.value = false
    }
}

function resetPage() {
    page.value = 1
    loadData()
}

function toggleAutoRefresh() {
    autoRefresh.value = !autoRefresh.value
    if (autoRefresh.value) {
        autoRefreshTimer = setInterval(loadData, 30000)
    } else if (autoRefreshTimer) {
        clearInterval(autoRefreshTimer)
        autoRefreshTimer = null
    }
}

function openDetail(event: EventLog) {
    selectedEvent.value = event
}

function closeDetail() {
    selectedEvent.value = null
}

function totalPages(): number {
    if (!response.value) return 1
    return Math.max(1, Math.ceil(response.value.total / response.value.per_page))
}

function handleFromChange(e: Event) {
    const val = (e.target as HTMLInputElement).value
    filterFrom.value = val ? fromLocalDatetimeValue(val) : ''
    resetPage()
}

function handleToChange(e: Event) {
    const val = (e.target as HTMLInputElement).value
    filterTo.value = val ? fromLocalDatetimeValue(val) : ''
    resetPage()
}

// Watch dropdown filters for immediate reload
watch([filterCategory, filterSource, filterEventType, filterEntityType], () => {
    resetPage()
})

// Debounce text search
watch(filterQuery, () => {
    if (debounceTimer) clearTimeout(debounceTimer)
    debounceTimer = setTimeout(() => {
        resetPage()
    }, 300)
})

onMounted(() => {
    // Default to last 24 hours
    const now = new Date()
    const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)
    filterFrom.value = yesterday.toISOString()
    loadData()
})

onUnmounted(() => {
    if (autoRefreshTimer) {
        clearInterval(autoRefreshTimer)
    }
    if (debounceTimer) {
        clearTimeout(debounceTimer)
    }
})
</script>

<template>
    <div class="mt-4">
        <h2>Activity</h2>
        <p class="text-muted mb-4">System event log</p>

        <!-- Filter bar -->
        <div class="d-flex gap-2 mb-3 flex-wrap align-items-end">
            <div>
                <label class="form-label form-label-sm mb-1">Category</label>
                <select v-model="filterCategory" class="form-select form-select-sm" style="min-width: 140px;">
                    <option value="">All</option>
                    <option v-for="c in response?.facets?.categories ?? []" :key="c" :value="c">{{ c }}</option>
                </select>
            </div>
            <div>
                <label class="form-label form-label-sm mb-1">Source</label>
                <select v-model="filterSource" class="form-select form-select-sm" style="min-width: 140px;">
                    <option value="">All</option>
                    <option v-for="s in response?.facets?.sources ?? []" :key="s" :value="s">{{ s }}</option>
                </select>
            </div>
            <div>
                <label class="form-label form-label-sm mb-1">Event Type</label>
                <select v-model="filterEventType" class="form-select form-select-sm" style="min-width: 160px;">
                    <option value="">All</option>
                    <option v-for="t in response?.facets?.event_types ?? []" :key="t" :value="t">{{ t }}</option>
                </select>
            </div>
            <div>
                <label class="form-label form-label-sm mb-1">Entity Type</label>
                <select v-model="filterEntityType" class="form-select form-select-sm" style="min-width: 140px;">
                    <option value="">All</option>
                    <option v-for="t in response?.facets?.entity_types ?? []" :key="t" :value="t">{{ t }}</option>
                </select>
            </div>
            <div>
                <label class="form-label form-label-sm mb-1">Search</label>
                <input v-model="filterQuery" type="text" class="form-control form-control-sm"
                    placeholder="Search summary..." style="min-width: 160px;">
            </div>
            <div>
                <label class="form-label form-label-sm mb-1">From</label>
                <input type="datetime-local" class="form-control form-control-sm" :value="toLocalDatetimeValue(filterFrom)"
                    @change="handleFromChange" style="min-width: 180px;">
            </div>
            <div>
                <label class="form-label form-label-sm mb-1">To</label>
                <input type="datetime-local" class="form-control form-control-sm" :value="toLocalDatetimeValue(filterTo)"
                    @change="handleToChange" style="min-width: 180px;">
            </div>
        </div>

        <!-- Toolbar -->
        <div class="d-flex justify-content-between align-items-center mb-3">
            <span class="text-muted" v-if="response">{{ response.total }} events</span>
            <span v-else>&nbsp;</span>
            <div class="d-flex gap-2">
                <button class="btn btn-sm" :class="autoRefresh ? 'btn-success' : 'btn-outline-secondary'"
                    @click="toggleAutoRefresh">
                    <i class="bi bi-arrow-repeat me-1"></i>
                    Auto-refresh {{ autoRefresh ? 'ON' : 'OFF' }}
                </button>
                <button class="btn btn-sm btn-outline-primary" @click="loadData" :disabled="loading">
                    <i class="bi bi-arrow-clockwise me-1"></i>
                    Refresh
                </button>
            </div>
        </div>

        <!-- Table -->
        <div class="position-relative" style="min-height: 200px;">
            <LoadingSpinner v-if="loading" />
            <table class="table table-hover table-sm" v-if="response">
                <thead>
                    <tr>
                        <th style="width: 100px;">Time</th>
                        <th style="width: 120px;">Category</th>
                        <th style="width: 180px;">Event</th>
                        <th style="width: 140px;">Source</th>
                        <th>Summary</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="event in response.items" :key="event.id" @click="openDetail(event)"
                        style="cursor: pointer;">
                        <td class="text-nowrap">
                            <small :title="new Date(event.created_at).toLocaleString()">{{
                                formatRelativeTime(event.created_at) }}</small>
                        </td>
                        <td>
                            <span class="badge" :class="categoryBadgeClass(event.category)">{{ event.category }}</span>
                        </td>
                        <td><small>{{ event.event_type }}</small></td>
                        <td>
                            <span class="badge" :class="sourceBadgeClass(event.source)">{{ event.source }}</span>
                        </td>
                        <td><small>{{ event.summary }}</small></td>
                    </tr>
                    <tr v-if="response.items.length === 0">
                        <td colspan="5" class="text-center text-muted py-4">No events found</td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Pagination -->
        <nav v-if="response && totalPages() > 1" class="d-flex justify-content-between align-items-center">
            <div>
                <select v-model.number="perPage" class="form-select form-select-sm" style="width: auto;"
                    @change="resetPage()">
                    <option :value="25">25 per page</option>
                    <option :value="50">50 per page</option>
                    <option :value="100">100 per page</option>
                    <option :value="200">200 per page</option>
                </select>
            </div>
            <ul class="pagination pagination-sm mb-0">
                <li class="page-item" :class="{ disabled: page <= 1 }">
                    <a class="page-link" href="#" @click.prevent="page--; loadData()">Prev</a>
                </li>
                <li class="page-item disabled">
                    <span class="page-link">Page {{ page }} of {{ totalPages() }}</span>
                </li>
                <li class="page-item" :class="{ disabled: page >= totalPages() }">
                    <a class="page-link" href="#" @click.prevent="page++; loadData()">Next</a>
                </li>
            </ul>
        </nav>

        <!-- Detail drawer -->
        <EventDetailDrawer :event="selectedEvent" @close="closeDetail" />
    </div>
</template>
