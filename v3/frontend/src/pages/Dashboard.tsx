import React from 'react'
import { useQuery } from 'react-query'
import {
  PageSection,
  Title,
  Card,
  CardBody,
  Grid,
  GridItem,
  Spinner,
  Alert,
  Flex,
  FlexItem,
} from '@patternfly/react-core'
import {
  ClusterIcon,
  ServerIcon,
  ShieldAltIcon,
} from '@patternfly/react-icons'
import { hubsAPI } from '@/services/api'

const Dashboard: React.FC = () => {
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
        <Alert variant="danger" title="Failed to load dashboard data">
          {(error as Error).message}
        </Alert>
      </div>
    )
  }

  const totalHubs = hubs?.length || 0
  const totalClusters =
    hubs?.reduce((acc, hub) => acc + (hub.managedClusters?.length || 0), 0) || 0
  const totalPolicies =
    hubs?.reduce((acc, hub) => acc + (hub.policiesInfo?.length || 0), 0) || 0
  const healthyHubs = hubs?.filter((hub) => hub.status === 'Ready').length || 0

  return (
    <>
      <PageSection>
        <Title headingLevel="h1" size="2xl">
          Dashboard
        </Title>
      </PageSection>

      <PageSection>
        <Grid hasGutter>
          <GridItem sm={12} md={6} lg={3}>
            <Card>
              <CardBody>
                <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  <FlexItem>
                    <Flex alignItems={{ default: 'alignItemsCenter' }}>
                      <FlexItem>
                        <ClusterIcon color="#06c" />
                      </FlexItem>
                      <FlexItem>
                        <Title headingLevel="h3">
                          Total Hubs
                        </Title>
                      </FlexItem>
                    </Flex>
                  </FlexItem>
                  <FlexItem>
                    <div className="value" style={{ fontSize: '2rem', fontWeight: 700 }}>
                      {totalHubs}
                    </div>
                  </FlexItem>
                  <FlexItem>
                    <small>
                      {healthyHubs} healthy / {totalHubs - healthyHubs} unhealthy
                    </small>
                  </FlexItem>
                </Flex>
              </CardBody>
            </Card>
          </GridItem>

          <GridItem sm={12} md={6} lg={3}>
            <Card>
              <CardBody>
                <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  <FlexItem>
                    <Flex alignItems={{ default: 'alignItemsCenter' }}>
                      <FlexItem>
                        <ServerIcon color="#009596" />
                      </FlexItem>
                      <FlexItem>
                        <Title headingLevel="h3">
                          Total Clusters
                        </Title>
                      </FlexItem>
                    </Flex>
                  </FlexItem>
                  <FlexItem>
                    <div className="value" style={{ fontSize: '2rem', fontWeight: 700 }}>
                      {totalClusters}
                    </div>
                  </FlexItem>
                  <FlexItem>
                    <small>Across all managed hubs</small>
                  </FlexItem>
                </Flex>
              </CardBody>
            </Card>
          </GridItem>

          <GridItem sm={12} md={6} lg={3}>
            <Card>
              <CardBody>
                <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  <FlexItem>
                    <Flex alignItems={{ default: 'alignItemsCenter' }}>
                      <FlexItem>
                        <ShieldAltIcon color="#8b2c9b" />
                      </FlexItem>
                      <FlexItem>
                        <Title headingLevel="h3">
                          Total Policies
                        </Title>
                      </FlexItem>
                    </Flex>
                  </FlexItem>
                  <FlexItem>
                    <div className="value" style={{ fontSize: '2rem', fontWeight: 700 }}>
                      {totalPolicies}
                    </div>
                  </FlexItem>
                  <FlexItem>
                    <small>Across all hubs</small>
                  </FlexItem>
                </Flex>
              </CardBody>
            </Card>
          </GridItem>

          <GridItem sm={12} md={6} lg={3}>
            <Card>
              <CardBody>
                <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  <FlexItem>
                    <Flex alignItems={{ default: 'alignItemsCenter' }}>
                      <FlexItem>
                        <ClusterIcon color="#f0ab00" />
                      </FlexItem>
                      <FlexItem>
                        <Title headingLevel="h3">
                          Health Rate
                        </Title>
                      </FlexItem>
                    </Flex>
                  </FlexItem>
                  <FlexItem>
                    <div className="value" style={{ fontSize: '2rem', fontWeight: 700 }}>
                      {totalHubs > 0 ? Math.round((healthyHubs / totalHubs) * 100) : 0}%
                    </div>
                  </FlexItem>
                  <FlexItem>
                    <small>Hub health percentage</small>
                  </FlexItem>
                </Flex>
              </CardBody>
            </Card>
          </GridItem>
        </Grid>
      </PageSection>

      <PageSection>
        <Title headingLevel="h2" size="xl">
          Recent Activity
        </Title>
        {hubs && hubs.length > 0 ? (
          <Card>
            <CardBody>
              <p>No recent activity to display.</p>
            </CardBody>
          </Card>
        ) : (
          <div className="empty-state">
            <p>No hubs found. Get started by connecting your first hub.</p>
          </div>
        )}
      </PageSection>
    </>
  )
}

export default Dashboard

