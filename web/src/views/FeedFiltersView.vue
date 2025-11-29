<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import FeedFilterRow from '@/components/FeedFilterRow.vue'
import type { FeedFilter, Feed, TorrentCategory, Notifier } from '@/types/api'

const toast = useToastStore()
const data = ref<FeedFilter[]>([])
const feeds = ref<Feed[]>([])
const categories = ref<TorrentCategory[]>([])
const notifiers = ref<Notifier[]>([])
const loading = ref(true)

const adding = ref(false)
const newFilter = ref<Partial<FeedFilter>>({})

onMounted(async () => {
  try {
    const [feedsData, categoriesData, notifiersData, filtersData] = await Promise.all([
      api.feeds.list(),
      api.torrentCategories.list(),
      api.notifiers.list(),
      api.feedFilters.list()
    ])
    feeds.value = feedsData
    categories.value = categoriesData
    notifiers.value = notifiersData
    data.value = filtersData
  } catch (e) {
    toast.error('Failed to load data')
  } finally {
    loading.value = false
  }
})

function startAdd() {
  adding.value = true
  newFilter.value = {
    name: '',
    feed_id: feeds.value[0]?.id,
    category_id: categories.value[0]?.id,
    notifier_id: notifiers.value[0]?.id
  }
}

function cancelAdd() {
  adding.value = false
  newFilter.value = {}
}

async function saveNew() {
  if (!newFilter.value.name || !newFilter.value.feed_id || !newFilter.value.category_id || !newFilter.value.notifier_id) {
    toast.error('All fields are required')
    return
  }

  try {
    const created = await api.feedFilters.create({
      name: newFilter.value.name,
      feed_id: Number(newFilter.value.feed_id),
      category_id: Number(newFilter.value.category_id),
      notifier_id: Number(newFilter.value.notifier_id)
    })
    data.value.push(created)
    toast.success('Feed filter created')
    adding.value = false
    newFilter.value = {}
  } catch (e) {
    toast.error('Failed to create feed filter')
  }
}

async function handleSave(filter: FeedFilter) {
  try {
    const updated = await api.feedFilters.update(filter.id, {
      name: filter.name,
      feed_id: Number(filter.feed_id),
      category_id: Number(filter.category_id),
      notifier_id: Number(filter.notifier_id)
    })
    const index = data.value.findIndex(d => d.id === filter.id)
    data.value[index] = updated
    toast.success('Feed filter updated')
  } catch (e) {
    toast.error('Failed to update feed filter')
  }
}

async function handleDelete(id: number) {
  try {
    await api.feedFilters.delete(id)
    data.value = data.value.filter(d => d.id !== id)
    toast.success('Feed filter deleted')
  } catch (e) {
    toast.error('Failed to delete feed filter')
  }
}
</script>

<template>
  <div class="mt-4">
    <h2>Feed Filters</h2>
    <p class="text-muted mb-4">Manage feed filters with expandable filter sets</p>

    <div class="position-relative">
      <LoadingSpinner v-if="loading" />

      <div class="mb-3">
        <button
          class="btn btn-primary btn-sm"
          @click="startAdd"
          :disabled="adding"
        >
          <i class="bi bi-plus-lg me-1"></i>
          Add New
        </button>
      </div>

      <table class="table table-dark table-striped table-hover">
        <thead>
          <tr>
            <th style="width: 40px"></th>
            <th>ID</th>
            <th>Name</th>
            <th>Feed</th>
            <th>Category</th>
            <th>Notifier</th>
            <th style="width: 120px">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="adding">
            <td></td>
            <td>-</td>
            <td>
              <input
                type="text"
                v-model="newFilter.name"
                class="form-control form-control-sm"
                placeholder="Filter name"
              />
            </td>
            <td>
              <select v-model="newFilter.feed_id" class="form-select form-select-sm">
                <option v-for="feed in feeds" :key="feed.id" :value="feed.id">
                  {{ feed.name }}
                </option>
              </select>
            </td>
            <td>
              <select v-model="newFilter.category_id" class="form-select form-select-sm">
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                  {{ cat.name }}
                </option>
              </select>
            </td>
            <td>
              <select v-model="newFilter.notifier_id" class="form-select form-select-sm">
                <option v-for="notifier in notifiers" :key="notifier.id" :value="notifier.id">
                  {{ notifier.name }}
                </option>
              </select>
            </td>
            <td>
              <button
                class="btn btn-success btn-sm me-1"
                @click="saveNew"
                title="Save"
              >
                <i class="bi bi-check"></i>
              </button>
              <button
                class="btn btn-secondary btn-sm"
                @click="cancelAdd"
                title="Cancel"
              >
                <i class="bi bi-x"></i>
              </button>
            </td>
          </tr>
          <tr v-if="data.length === 0 && !adding">
            <td colspan="7" class="text-center text-muted py-4">
              No feed filters available
            </td>
          </tr>
          <FeedFilterRow
            v-for="filter in data"
            :key="filter.id"
            :filter="filter"
            :feeds="feeds"
            :categories="categories"
            :notifiers="notifiers"
            @save="handleSave"
            @delete="handleDelete"
          />
        </tbody>
      </table>
    </div>
  </div>
</template>
