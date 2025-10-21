import { formatDistanceToNow, format } from 'date-fns'

export const formatDate = (dateString: string): string => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return format(date, 'yyyy-MM-dd HH:mm:ss')
}

export const formatDateRelative = (dateString: string): string => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return formatDistanceToNow(date, { addSuffix: true })
}

export const getStatusColor = (status: string): 'green' | 'red' | 'orange' | 'blue' => {
  const statusLower = status.toLowerCase()
  if (statusLower === 'ready' || statusLower === 'available' || statusLower === 'healthy') {
    return 'green'
  }
  if (statusLower === 'notready' || statusLower === 'unavailable' || statusLower === 'failed') {
    return 'red'
  }
  if (statusLower === 'pending' || statusLower === 'progressing') {
    return 'orange'
  }
  return 'blue'
}

export const getComplianceColor = (state: string): 'green' | 'red' | 'orange' => {
  const stateLower = state.toLowerCase()
  if (stateLower === 'compliant') return 'green'
  if (stateLower === 'noncompliant') return 'red'
  return 'orange'
}

export const getSeverityColor = (severity: string): 'green' | 'red' | 'orange' | 'blue' => {
  const severityLower = severity.toLowerCase()
  if (severityLower === 'critical' || severityLower === 'high') return 'red'
  if (severityLower === 'medium') return 'orange'
  if (severityLower === 'low') return 'blue'
  return 'green'
}

export const formatBytes = (bytes: string): string => {
  if (!bytes) return 'N/A'
  
  // Parse the byte string (e.g., "1024Ki", "1Gi", etc.)
  const units = ['Ki', 'Mi', 'Gi', 'Ti', 'Pi']
  
  for (const unit of units) {
    if (bytes.endsWith(unit)) {
      const value = parseFloat(bytes.slice(0, -unit.length))
      return `${value.toFixed(2)} ${unit}B`
    }
  }
  
  // If no unit found, return as is
  return bytes
}

export const formatCPU = (cpu: string): string => {
  if (!cpu) return 'N/A'
  
  // CPU can be in cores (e.g., "4") or millicores (e.g., "4000m")
  if (cpu.endsWith('m')) {
    const millicores = parseInt(cpu.slice(0, -1))
    return `${(millicores / 1000).toFixed(2)} cores`
  }
  
  return `${cpu} cores`
}

