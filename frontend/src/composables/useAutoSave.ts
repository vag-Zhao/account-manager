import { ref } from 'vue'

/**
 * Composable for debounced auto-save functionality
 * @param saveCallback - The function to call when saving
 * @param delay - Debounce delay in milliseconds (default: 300)
 */
export function useAutoSave(saveCallback: () => Promise<void>, delay = 300) {
  const saveTimer = ref<number | null>(null)
  const isSaving = ref(false)

  const triggerAutoSave = () => {
    if (saveTimer.value) {
      clearTimeout(saveTimer.value)
    }
    saveTimer.value = window.setTimeout(async () => {
      isSaving.value = true
      try {
        await saveCallback()
      } finally {
        isSaving.value = false
      }
    }, delay)
  }

  const cancelAutoSave = () => {
    if (saveTimer.value) {
      clearTimeout(saveTimer.value)
      saveTimer.value = null
    }
  }

  return {
    triggerAutoSave,
    cancelAutoSave,
    isSaving
  }
}
