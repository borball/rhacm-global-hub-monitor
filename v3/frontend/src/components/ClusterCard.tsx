import React from 'react'
import {
  Card,
  CardTitle,
  CardBody,
  CardFooter,
  Label,
  Flex,
  FlexItem,
  Title,
} from '@patternfly/react-core'
import { ClusterIcon } from '@patternfly/react-icons'
import type { ManagedCluster, ManagedHub } from '@/types'
import StatusLabel from './StatusLabel'
import { formatDateRelative } from '@/utils/helpers'

interface ClusterCardProps {
  cluster: ManagedCluster | ManagedHub
  hubName?: string
  isHub?: boolean
}

const ClusterCard: React.FC<ClusterCardProps> = ({ cluster, hubName, isHub = false }) => {
  const clusterCount =
    isHub && 'managedClusters' in cluster ? cluster.managedClusters.length : 0
  const linkTo = isHub
    ? `/hubs/${cluster.name}`
    : `/hubs/${hubName}/clusters/${cluster.name}`

  return (
    <Card isSelectable isClickable component={'div' as any} onClick={() => window.location.href = linkTo} style={{ textDecoration: 'none', cursor: 'pointer' }}>
      <CardTitle>
        <Flex alignItems={{ default: 'alignItemsCenter' }}>
          <FlexItem>
            <ClusterIcon />
          </FlexItem>
          <FlexItem>
            <Title headingLevel="h3" size="lg">
              {cluster.name}
            </Title>
          </FlexItem>
        </Flex>
      </CardTitle>
      <CardBody>
        <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
          <FlexItem>
            <strong>Status:</strong> <StatusLabel status={cluster.status} />
          </FlexItem>
          <FlexItem>
            <strong>Version:</strong> {cluster.version || 'N/A'}
          </FlexItem>
          <FlexItem>
            <strong>Platform:</strong> {cluster.clusterInfo.platform || 'N/A'}
          </FlexItem>
          {isHub && (
            <FlexItem>
              <strong>Managed Clusters:</strong>{' '}
              <Label color="blue">{clusterCount}</Label>
            </FlexItem>
          )}
          <FlexItem>
            <strong>Nodes:</strong> <Label color="cyan">{cluster.nodesInfo.length}</Label>
          </FlexItem>
          <FlexItem>
            <strong>Policies:</strong> <Label color="purple">{cluster.policiesInfo.length}</Label>
          </FlexItem>
        </Flex>
      </CardBody>
      <CardFooter>
        <small>Created {formatDateRelative(cluster.createdAt)}</small>
      </CardFooter>
    </Card>
  )
}

export default ClusterCard

