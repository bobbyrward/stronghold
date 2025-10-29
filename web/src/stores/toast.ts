import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Toast {
  id: number
  message: string
  type: 'success' | 'error' | 'warning' | 'info'
  duration: number
}

export const useToastStore = defineStore('toast', () => {
  const toasts = ref<Toast[]>([])
  let nextId = 0

  function add(message: string, type: Toast['type'], duration = 5000) {
    const id = nextId++
    toasts.value.push({ id, message, type, duration })
    setTimeout(() => remove(id), duration)
  }

  function remove(id: number) {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  const success = (msg: string) => add(msg, 'success')
  const error = (msg: string) => add(msg, 'error')
  const warning = (msg: string) => add(msg, 'warning')
  const info = (msg: string) => add(msg, 'info')

  return { toasts, add, remove, success, error, warning, info }
})
