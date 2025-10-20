import React from 'react'
import { Label } from '@patternfly/react-core'
import {
  CheckCircleIcon,
  ExclamationCircleIcon,
  ExclamationTriangleIcon,
  InfoCircleIcon,
} from '@patternfly/react-icons'
import { getStatusColor } from '@/utils/helpers'

interface StatusLabelProps {
  status: string
  icon?: boolean
}

const StatusLabel: React.FC<StatusLabelProps> = ({ status, icon = true }) => {
  const color = getStatusColor(status)

  const getIcon = () => {
    switch (color) {
      case 'green':
        return <CheckCircleIcon />
      case 'red':
        return <ExclamationCircleIcon />
      case 'orange':
        return <ExclamationTriangleIcon />
      default:
        return <InfoCircleIcon />
    }
  }

  return (
    <Label color={color} icon={icon ? getIcon() : undefined}>
      {status}
    </Label>
  )
}

export default StatusLabel

