/**
 * Composable for date formatting utilities
 * Uses native Date API to avoid external dependencies
 */
export function useDateFormat() {
  const formatDate = (date: string | Date | null | undefined, format = 'YYYY-MM-DD'): string => {
    if (!date) return ''
    const d = new Date(date)
    if (isNaN(d.getTime())) return ''

    const year = d.getFullYear()
    const month = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    const hours = String(d.getHours()).padStart(2, '0')
    const minutes = String(d.getMinutes()).padStart(2, '0')
    const seconds = String(d.getSeconds()).padStart(2, '0')

    return format
      .replace('YYYY', String(year))
      .replace('MM', month)
      .replace('DD', day)
      .replace('HH', hours)
      .replace('mm', minutes)
      .replace('ss', seconds)
  }

  const formatDateTime = (date: string | Date | null | undefined): string => {
    return formatDate(date, 'YYYY-MM-DD HH:mm:ss')
  }

  const formatRelative = (date: string | Date | null | undefined): string => {
    if (!date) return ''
    const now = new Date()
    const target = new Date(date)
    if (isNaN(target.getTime())) return ''

    const diffTime = now.getTime() - target.getTime()
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))

    if (diffDays === 0) return '今天'
    if (diffDays === 1) return '昨天'
    if (diffDays === -1) return '明天'
    if (diffDays > 0 && diffDays < 7) return `${diffDays}天前`
    if (diffDays < 0 && diffDays > -7) return `${Math.abs(diffDays)}天后`

    return formatDate(date)
  }

  const isExpired = (date: string | Date | null | undefined): boolean => {
    if (!date) return false
    const target = new Date(date)
    if (isNaN(target.getTime())) return false
    const now = new Date()
    now.setHours(0, 0, 0, 0)
    target.setHours(0, 0, 0, 0)
    return target < now
  }

  const isExpiringSoon = (date: string | Date | null | undefined, days = 7): boolean => {
    if (!date) return false
    const target = new Date(date)
    if (isNaN(target.getTime())) return false
    const now = new Date()
    const diffTime = target.getTime() - now.getTime()
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    return diffDays >= 0 && diffDays <= days
  }

  return {
    formatDate,
    formatDateTime,
    formatRelative,
    isExpired,
    isExpiringSoon
  }
}
