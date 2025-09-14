import React from 'react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Sidebar from './components/Sidebar/Sidebar'
import Main from './components/Main/Main'
import Home from './components/Home/Home'
import Loading from './components/Loading/Loading'

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/loading" element={<Loading />} />
        <Route path="/main" element={
          <>
            <Sidebar />
            <Main />
          </>
        } />
      </Routes>
    </Router>
  )
}

export default App