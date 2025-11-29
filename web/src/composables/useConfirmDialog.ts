import { ref } from 'vue'

export function useConfirmDialog() {
  const show = ref(false)
  const title = ref('')
  const message = ref('')
  const resolvePromise = ref<((value: boolean) => void) | null>(null)

  function confirm(dialogTitle: string, dialogMessage: string): Promise<boolean> {
    title.value = dialogTitle
    message.value = dialogMessage
    show.value = true

    return new Promise<boolean>((resolve) => {
      resolvePromise.value = resolve
    })
  }

  function handleConfirm() {
    show.value = false
    if (resolvePromise.value) {
      resolvePromise.value(true)
      resolvePromise.value = null
    }
  }

  function handleCancel() {
    show.value = false
    if (resolvePromise.value) {
      resolvePromise.value(false)
      resolvePromise.value = null
    }
  }

  return {
    show,
    title,
    message,
    confirm,
    handleConfirm,
    handleCancel
  }
}
