import React from 'react'
import {
  Card,
  CardTitle,
  CardBody,
  DescriptionList,
  DescriptionListGroup,
  DescriptionListTerm,
  DescriptionListDescription,
  Label,
  Flex,
  FlexItem,
} from '@patternfly/react-core'
import { ServerIcon } from '@patternfly/react-icons'
import type { NodeInfo } from '@/types'
import StatusLabel from './StatusLabel'

interface NodeCardProps {
  node: NodeInfo
}

const NodeCard: React.FC<NodeCardProps> = ({ node }) => {
  return (
    <Card isCompact>
      <CardTitle>
        <Flex alignItems={{ default: 'alignItemsCenter' }}>
          <FlexItem>
            <ServerIcon />
          </FlexItem>
          <FlexItem>
            <strong>{node.name}</strong>
          </FlexItem>
          <FlexItem>
            <StatusLabel status={node.status} />
          </FlexItem>
        </Flex>
      </CardTitle>
      <CardBody>
        <DescriptionList isCompact isHorizontal columnModifier={{ default: '2Col' }}>
          {node.capacity.cpu && (
            <DescriptionListGroup>
              <DescriptionListTerm>CPU</DescriptionListTerm>
              <DescriptionListDescription>
                <Label color="blue">{node.capacity.cpu}</Label>
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {node.capacity.memory && (
            <DescriptionListGroup>
              <DescriptionListTerm>RAM</DescriptionListTerm>
              <DescriptionListDescription>
                <Label color="green">{node.capacity.memory}</Label>
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {node.capacity.storage && (
            <DescriptionListGroup>
              <DescriptionListTerm>Storage</DescriptionListTerm>
              <DescriptionListDescription>
                <Label color="orange">{node.capacity.storage}</Label>
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {node.internalIP && (
            <DescriptionListGroup>
              <DescriptionListTerm>IP Address</DescriptionListTerm>
              <DescriptionListDescription>{node.internalIP}</DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {node.annotations?.['bmc-address'] && (
            <DescriptionListGroup>
              <DescriptionListTerm>BMC</DescriptionListTerm>
              <DescriptionListDescription>
                <small style={{ fontFamily: 'monospace' }}>
                  {node.annotations['bmc-address']}
                </small>
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {node.annotations?.manufacturer && (
            <DescriptionListGroup>
              <DescriptionListTerm>Manufacturer</DescriptionListTerm>
              <DescriptionListDescription>
                {node.annotations.manufacturer}
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {node.annotations?.['product-name'] && (
            <DescriptionListGroup>
              <DescriptionListTerm>Product</DescriptionListTerm>
              <DescriptionListDescription>
                {node.annotations['product-name']}
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {node.annotations?.['serial-number'] && (
            <DescriptionListGroup>
              <DescriptionListTerm>Serial Number</DescriptionListTerm>
              <DescriptionListDescription>
                <small style={{ fontFamily: 'monospace' }}>
                  {node.annotations['serial-number']}
                </small>
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
          {node.annotations?.['nic-count'] && (
            <DescriptionListGroup>
              <DescriptionListTerm>NICs</DescriptionListTerm>
              <DescriptionListDescription>
                {node.annotations['nic-count']} interfaces
              </DescriptionListDescription>
            </DescriptionListGroup>
          )}
        </DescriptionList>
      </CardBody>
    </Card>
  )
}

export default NodeCard

