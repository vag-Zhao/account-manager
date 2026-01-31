/**
 * Generic debounce composable
 * @param fn - Function to debounce
 * @param delay - Delay in milliseconds
 */
export function useDebounce<T extends (...args: any[]) => any>(
  fn: T,
  delay = 300
): (...args: Parameters<T>) => void {
  let timeoutId: number | null = null

  return (...args: Parameters<T>) => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }
    timeoutId = window.setTimeout(() => {
      fn(...args)
    }, delay)
  }
}
