import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import {
  Page,
  Masthead,
  MastheadToggle,
  MastheadMain,
  MastheadBrand,
  MastheadContent,
  PageSidebar,
  PageSidebarBody,
  Nav,
  NavList,
  NavItem,
  PageSection,
  Title,
} from '@patternfly/react-core'
import { BarsIcon } from '@patternfly/react-icons'

interface LayoutProps {
  children: React.ReactNode
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [isSidebarOpen, setIsSidebarOpen] = React.useState(true)
  const location = useLocation()

  const Header = (
    <Masthead>
      <MastheadToggle>
        <button
          type="button"
          onClick={() => setIsSidebarOpen(!isSidebarOpen)}
          style={{
            background: 'none',
            border: 'none',
            cursor: 'pointer',
            padding: '8px',
            display: 'flex',
            alignItems: 'center',
          }}
        >
          <BarsIcon />
        </button>
      </MastheadToggle>
      <MastheadMain>
        <MastheadBrand>
          <Title headingLevel="h1" size="2xl">
            RHACM Global Hub Monitor
          </Title>
        </MastheadBrand>
      </MastheadMain>
      <MastheadContent>{/* Add user menu here if needed */}</MastheadContent>
    </Masthead>
  )

  const Sidebar = (
    <PageSidebar isSidebarOpen={isSidebarOpen}>
      <PageSidebarBody>
        <Nav>
          <NavList>
            <NavItem
              isActive={location.pathname === '/dashboard'}
              to="/dashboard"
              component={(props: any) => <Link {...props} />}
            >
              Dashboard
            </NavItem>
            <NavItem
              isActive={location.pathname.startsWith('/hubs')}
              to="/hubs"
              component={(props: any) => <Link {...props} />}
            >
              Managed Hubs
            </NavItem>
          </NavList>
        </Nav>
      </PageSidebarBody>
    </PageSidebar>
  )

  return (
    <Page header={Header} sidebar={Sidebar}>
      <PageSection>{children}</PageSection>
    </Page>
  )
}

export default Layout

