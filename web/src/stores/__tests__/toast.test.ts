import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useToastStore } from '../toast'

describe('toast store', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('add() adds toast to list', () => {
    const store = useToastStore()

    store.add('Test message', 'info')

    expect(store.toasts.length).toBe(1)
    expect(store.toasts[0]!.message).toBe('Test message')
    expect(store.toasts[0]!.type).toBe('info')
  })

  it('remove() removes toast by id', () => {
    const store = useToastStore()

    store.add('Toast 1', 'info')
    store.add('Toast 2', 'info')

    const toastId = store.toasts[0]!.id
    store.remove(toastId)

    expect(store.toasts.length).toBe(1)
    expect(store.toasts[0]!.message).toBe('Toast 2')
  })

  it('success() convenience method works', () => {
    const store = useToastStore()

    store.success('Success message')

    expect(store.toasts.length).toBe(1)
    expect(store.toasts[0]!.message).toBe('Success message')
    expect(store.toasts[0]!.type).toBe('success')
  })

  it('error() convenience method works', () => {
    const store = useToastStore()

    store.error('Error message')

    expect(store.toasts.length).toBe(1)
    expect(store.toasts[0]!.message).toBe('Error message')
    expect(store.toasts[0]!.type).toBe('error')
  })

  it('warning() convenience method works', () => {
    const store = useToastStore()

    store.warning('Warning message')

    expect(store.toasts.length).toBe(1)
    expect(store.toasts[0]!.message).toBe('Warning message')
    expect(store.toasts[0]!.type).toBe('warning')
  })

  it('info() convenience method works', () => {
    const store = useToastStore()

    store.info('Info message')

    expect(store.toasts.length).toBe(1)
    expect(store.toasts[0]!.message).toBe('Info message')
    expect(store.toasts[0]!.type).toBe('info')
  })

  it('auto-dismiss removes toast after duration', () => {
    const store = useToastStore()

    store.add('Auto dismiss', 'info', 3000)

    expect(store.toasts.length).toBe(1)

    vi.advanceTimersByTime(3000)

    expect(store.toasts.length).toBe(0)
  })

  it('generates unique ids for each toast', () => {
    const store = useToastStore()

    store.add('Toast 1', 'info')
    store.add('Toast 2', 'info')

    expect(store.toasts[0]!.id).not.toBe(store.toasts[1]!.id)
  })
})
