export function statusSeverity(status) {
  switch (status) {
    case 'Completed': return 'success'
    case 'Failed': return 'danger'
    case 'PartiallyFailed': return 'warn'
    case 'InProgress': return 'info'
    case 'Deleting': return 'secondary'
    default: return 'secondary'
  }
}

export function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString()
}

export function formatBytes(bytes) {
  if (bytes === 0 || bytes == null) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const idx = Math.min(i, units.length - 1)
  return (bytes / Math.pow(k, idx)).toFixed(1) + ' ' + units[idx]
}
