import React from 'react'
import ReactDOM from 'react-dom/client'
import RootLayout from './root-layout'
import { AuthProvider } from './api/auth-provider'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AuthProvider>
      <RootLayout />
    </AuthProvider>
  </React.StrictMode>,
)
