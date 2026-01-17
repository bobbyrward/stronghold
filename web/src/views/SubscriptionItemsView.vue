<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { Author, AuthorSubscription, AuthorSubscriptionItem } from '@/types/api'

interface SubscriptionItemWithAuthor extends AuthorSubscriptionItem {
    author_id: number
    author_name: string
}

const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const loading = ref(true)
const authors = ref<Author[]>([])
const authorsWithSubscriptions = ref<{ id: number; name: string }[]>([])
const allItems = ref<SubscriptionItemWithAuthor[]>([])
const selectedAuthorId = ref<number | null>(null)

// Filtered and sorted items
const filteredItems = computed(() => {
    let items = [...allItems.value]

    // Filter by author if selected
    if (selectedAuthorId.value !== null) {
        items = items.filter(item => item.author_id === selectedAuthorId.value)
    }

    // Sort by downloaded_at descending (newest first)
    items.sort((a, b) => new Date(b.downloaded_at).getTime() - new Date(a.downloaded_at).getTime())

    return items
})

onMounted(async () => {
    // Check for author_id query parameter
    const authorIdParam = route.query.author_id
    if (authorIdParam) {
        selectedAuthorId.value = parseInt(authorIdParam as string, 10)
    }

    await loadData()
})

// Update URL when filter changes
watch(selectedAuthorId, (newValue) => {
    const query = newValue !== null ? { author_id: String(newValue) } : {}
    router.replace({ query })
})

async function loadData() {
    loading.value = true

    try {
        // Load all authors
        authors.value = await api.authors.list()

        // For each author, try to get subscription and items
        const authorsWithSubs: { id: number; name: string }[] = []
        const items: SubscriptionItemWithAuthor[] = []

        for (const author of authors.value) {
            try {
                // Try to get subscription
                const subscription: AuthorSubscription = await api.authors.subscription.get(author.id)
                authorsWithSubs.push({ id: author.id, name: author.name })

                // Get subscription items
                try {
                    const subItems = await api.authors.subscription.items(author.id)
                    for (const item of subItems) {
                        items.push({
                            ...item,
                            author_id: author.id,
                            author_name: subscription.author_name
                        })
                    }
                } catch {
                    // No items for this subscription, that's fine
                }
            } catch {
                // No subscription for this author, skip
            }
        }

        authorsWithSubscriptions.value = authorsWithSubs
        allItems.value = items

        // Validate selected author still has a subscription
        if (selectedAuthorId.value !== null) {
            const hasSubscription = authorsWithSubs.some(a => a.id === selectedAuthorId.value)
            if (!hasSubscription) {
                selectedAuthorId.value = null
            }
        }
    } catch (e) {
        toast.error('Failed to load subscription items')
    } finally {
        loading.value = false
    }
}

function formatDate(dateString: string): string {
    const date = new Date(dateString)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}`
}

function truncateHash(hash: string): string {
    if (hash.length <= 8) return hash
    return hash.substring(0, 8) + '...'
}

function handleAuthorChange(event: Event) {
    const target = event.target as HTMLSelectElement
    const value = target.value
    selectedAuthorId.value = value === '' ? null : parseInt(value, 10)
}
</script>

<template>
    <div class="mt-4">
        <h2>Download History</h2>
        <p class="text-muted mb-4">View subscription download history</p>

        <div class="position-relative">
            <LoadingSpinner v-if="loading" />

            <!-- Filter controls -->
            <div class="row mb-3">
                <div class="col-auto">
                    <label for="authorFilter" class="col-form-label">Filter by Author:</label>
                </div>
                <div class="col-auto">
                    <select id="authorFilter" class="form-select form-select-sm" :value="selectedAuthorId ?? ''"
                        @change="handleAuthorChange">
                        <option value="">All Authors</option>
                        <option v-for="author in authorsWithSubscriptions" :key="author.id" :value="author.id">
                            {{ author.name }}
                        </option>
                    </select>
                </div>
            </div>

            <!-- Empty state -->
            <div v-if="!loading && filteredItems.length === 0" class="text-center text-muted py-4">
                <template v-if="allItems.length === 0">
                    No subscription downloads found.
                </template>
                <template v-else>
                    No downloads found for the selected author.
                </template>
            </div>

            <!-- Items table -->
            <table v-if="filteredItems.length > 0" class="table table-dark table-striped table-hover">
                <thead>
                    <tr>
                        <th>Author</th>
                        <th>Downloaded At</th>
                        <th>Booksearch ID</th>
                        <th>Torrent Hash</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="item in filteredItems" :key="item.id">
                        <td>{{ item.author_name }}</td>
                        <td>{{ formatDate(item.downloaded_at) }}</td>
                        <td>{{ item.booksearch_id }}</td>
                        <td :title="item.torrent_hash">{{ truncateHash(item.torrent_hash) }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>
