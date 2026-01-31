/**
 * Composable for copying text to clipboard with formatting
 */
export function useCopyToClipboard() {
  // Import message dynamically to avoid build errors
  const showMessage = (type: 'success' | 'error', text: string) => {
    // This will be handled by the component using this composable
    console.log(`[${type}] ${text}`)
  }
  const copyToClipboard = async (
    text: string,
    onSuccess?: () => void,
    onError?: () => void
  ) => {
    try {
      await navigator.clipboard.writeText(text)
      if (onSuccess) onSuccess()
    } catch (err) {
      // Fallback for older browsers
      const textArea = document.createElement('textarea')
      textArea.value = text
      textArea.style.position = 'fixed'
      textArea.style.left = '-999999px'
      document.body.appendChild(textArea)
      textArea.select()
      try {
        document.execCommand('copy')
        if (onSuccess) onSuccess()
      } catch (e) {
        if (onError) onError()
      }
      document.body.removeChild(textArea)
    }
  }

  const copyWithFormat = async (
    account: string,
    password: string,
    format: string,
    onSuccess?: () => void,
    onError?: () => void
  ) => {
    const text = format
      .replace(/\{account\}/g, account)
      .replace(/\{password\}/g, password)
    await copyToClipboard(text, onSuccess, onError)
  }

  return {
    copyToClipboard,
    copyWithFormat
  }
}
