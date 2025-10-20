import React from 'react'
import {
  Table,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
} from '@patternfly/react-table'
import { Label } from '@patternfly/react-core'
import type { PolicyInfo } from '@/types'
import { getComplianceColor } from '@/utils/helpers'

interface PolicyTableProps {
  policies: PolicyInfo[]
}

const PolicyTable: React.FC<PolicyTableProps> = ({ policies }) => {
  if (policies.length === 0) {
    return <p>No policies found</p>
  }

  return (
    <Table aria-label="Policies table" variant="compact">
      <Thead>
        <Tr>
          <Th>Policy Name</Th>
          <Th>Compliance</Th>
          <Th>Remediation</Th>
          <Th>Standards</Th>
        </Tr>
      </Thead>
      <Tbody>
        {policies.map((policy) => (
          <Tr key={policy.name}>
            <Td dataLabel="Policy Name">
              <div>
                <strong>{policy.name.split('.').pop()}</strong>
                <br />
                <small style={{ color: '#6a6e73' }}>{policy.namespace}</small>
              </div>
            </Td>
            <Td dataLabel="Compliance">
              <Label color={getComplianceColor(policy.complianceState)}>
                {policy.complianceState || 'Unknown'}
              </Label>
            </Td>
            <Td dataLabel="Remediation">
              <Label color={policy.remediationAction === 'enforce' ? 'orange' : 'blue'}>
                {policy.remediationAction || 'N/A'}
              </Label>
            </Td>
            <Td dataLabel="Standards">
              {policy.standards && policy.standards.length > 0 ? (
                policy.standards.map((std, idx) => (
                  <Label key={idx} color="purple" style={{ marginRight: '4px' }}>
                    {std}
                  </Label>
                ))
              ) : (
                'N/A'
              )}
            </Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  )
}

export default PolicyTable

