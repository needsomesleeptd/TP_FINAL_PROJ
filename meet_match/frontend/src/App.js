import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Session from './components/Session'
import CurSession from './components/CurSession'
import Cards from './components/Cards'
import NotFound from './components/NotFound'

function App() {
    return (
      <Router>
        <Routes>
          <Route path='/session' element={<Session />} />
          <Route path='/session/:id' element={<CurSession />} />
          <Route path='/cards' element={<Cards />} />
          <Route path='/' element={<Navigate replace to="/session" />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </Router>
    );
  }

export default App
