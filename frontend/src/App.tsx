import React from 'react'
import Sidebar from './components/Sidebar/Sidebar'
import Main from './components/Main/Main'

const App = () => {
  return (
    <>
    <Sidebar/>
    <Main/>
      
    </>
  )
}

export default App






















/**
 * 
 * From my initial loading one, keep cause will match later, especially the color
 * 
 * 
 * 
 * import './App.css'
import { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './components/Home'
import Loading from './components/Loading';

function App() {
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const timer = setTimeout(() => setLoading(false), 5000);
        return () => clearTimeout(timer);
    }, []);

    return (
        <Router>
            <Routes>
                <Route path="/" element={<Home/>} />
            </Routes>
        </Router>
    )
}

export default App;
 */