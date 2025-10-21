import React from 'react'
import { useQuery } from 'react-query'
import {
  PageSection,
  Title,
  Spinner,
  Alert,
  Grid,
  GridItem,
  EmptyState,
  EmptyStateIcon,
  EmptyStateBody,
} from '@patternfly/react-core'
import { ClusterIcon } from '@patternfly/react-icons'
import { hubsAPI } from '@/services/api'
import ClusterCard from '@/components/ClusterCard'

const HubsList: React.FC = () => {
  const { data: hubs, isLoading, error } = useQuery('hubs', hubsAPI.listHubs)

  if (isLoading) {
    return (
      <div className="loading-container">
        <Spinner size="xl" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="error-container">
        <Alert variant="danger" title="Failed to load hubs">
          {(error as Error).message}
        </Alert>
      </div>
    )
  }

  return (
    <>
      <PageSection>
        <Title headingLevel="h1" size="2xl">
          Managed Hubs
        </Title>
        <p style={{ marginTop: '0.5rem', color: '#6a6e73' }}>
          View and manage all RHACM hub clusters
        </p>
      </PageSection>

      <PageSection>
        {hubs && hubs.length > 0 ? (
          <Grid hasGutter>
            {hubs.map((hub) => (
              <GridItem key={hub.name} sm={12} md={6} lg={4}>
                <ClusterCard cluster={hub} isHub />
              </GridItem>
            ))}
          </Grid>
        ) : (
          <EmptyState>
            <EmptyStateIcon icon={ClusterIcon} />
            <Title headingLevel="h2" size="lg">
              No managed hubs found
            </Title>
            <EmptyStateBody>
              There are no managed hub clusters configured. Please ensure your RHACM global hub is
              properly configured and has managed hubs registered.
            </EmptyStateBody>
          </EmptyState>
        )}
      </PageSection>
    </>
  )
}

export default HubsList

