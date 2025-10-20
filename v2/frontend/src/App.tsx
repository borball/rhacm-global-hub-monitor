import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from 'react-query'
import '@patternfly/react-core/dist/styles/base.css'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import HubsList from './pages/HubsList'
import HubDetail from './pages/HubDetail'
import ClusterDetail from './pages/ClusterDetail'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 30000, // 30 seconds
    },
  },
})

const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Layout>
          <Routes>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/hubs" element={<HubsList />} />
            <Route path="/hubs/:hubName" element={<HubDetail />} />
            <Route path="/hubs/:hubName/clusters/:clusterName" element={<ClusterDetail />} />
          </Routes>
        </Layout>
      </Router>
    </QueryClientProvider>
  )
}

export default App

