// import { StrictMode } from 'react' // Temporarily disabled
import { createRoot } from 'react-dom/client'
import './index.css'
import './i18n/index.ts'
import App from './App.tsx'

const rootEl = typeof document !== 'undefined' ? document.getElementById('root') : null

if (rootEl) {
  createRoot(rootEl).render(
      <App />
  )
}
