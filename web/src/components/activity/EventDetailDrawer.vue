<script setup lang="ts">
import type { EventLog } from '@/types/api'

defineProps<{
    event: EventLog | null
}>()

const emit = defineEmits<{
    close: []
}>()

function formatDetails(details: string): string {
    if (!details) return ''
    try {
        return JSON.stringify(JSON.parse(details), null, 2)
    } catch {
        return details
    }
}

function formatTimestamp(iso: string): string {
    return new Date(iso).toLocaleString()
}

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
</script>

<template>
    <Teleport to="body">
        <div v-if="event" class="offcanvas offcanvas-end show" style="visibility: visible" tabindex="-1">
            <div class="offcanvas-header">
                <h5 class="offcanvas-title">Event Detail</h5>
                <button type="button" class="btn-close" aria-label="Close" @click="emit('close')"></button>
            </div>
            <div class="offcanvas-body">
                <dl class="row mb-0">
                    <dt class="col-sm-4">Time</dt>
                    <dd class="col-sm-8">{{ formatTimestamp(event.created_at) }}</dd>

                    <dt class="col-sm-4">Category</dt>
                    <dd class="col-sm-8">
                        <span class="badge" :class="categoryBadgeClass(event.category)">{{ event.category }}</span>
                    </dd>

                    <dt class="col-sm-4">Event Type</dt>
                    <dd class="col-sm-8">{{ event.event_type }}</dd>

                    <dt class="col-sm-4">Source</dt>
                    <dd class="col-sm-8">
                        <span class="badge" :class="sourceBadgeClass(event.source)">{{ event.source }}</span>
                    </dd>

                    <dt class="col-sm-4">Entity</dt>
                    <dd class="col-sm-8">
                        <span v-if="event.entity_type">{{ event.entity_type }}</span>
                        <span v-if="event.entity_type && event.entity_id"> / </span>
                        <span v-if="event.entity_id">{{ event.entity_id }}</span>
                        <span v-if="!event.entity_type && !event.entity_id" class="text-muted">-</span>
                    </dd>

                    <dt class="col-sm-4">Summary</dt>
                    <dd class="col-sm-8">{{ event.summary }}</dd>

                    <dt class="col-12 mt-2">Details</dt>
                    <dd class="col-12">
                        <pre v-if="event.details" class="bg-dark text-light p-3 rounded small mt-1"
                            style="max-height: 400px; overflow: auto; white-space: pre-wrap;">{{ formatDetails(event.details) }}</pre>
                        <span v-else class="text-muted">No details</span>
                    </dd>
                </dl>
            </div>
        </div>
        <div v-if="event" class="offcanvas-backdrop fade show" @click="emit('close')"></div>
    </Teleport>
</template>
