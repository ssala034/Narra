import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'
import ContextProvider from './Context/Context.jsx'

// note he used the context provider here so can work with the api?? might change for me based on how I use it for Golang & wails
// will be different cause I am using typescript also
ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <ContextProvider>
    <App />
  </ContextProvider>
)
