// import { StrictMode } from 'react' // Temporarily disabled
import { createRoot } from 'react-dom/client'
import './index.css'
import './i18n/index.ts'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
    <App />
)
