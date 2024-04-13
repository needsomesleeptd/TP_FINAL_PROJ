import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import Registration from './components/Registration'
import Login from './components/Login'
import Main from './components/Main'
import Session from './components/Session'
import Cards from './components/Cards'
import NotFound from './components/NotFound'
import './custom.css'

function App() {
  const [cookies] = useCookies(['AccessToken']);
  const isLoggedIn = !!cookies.AccessToken;
  const [showLogin, setShowLogin] = useState(true);

  const requireAuth = (element) => {
    return isLoggedIn ? element : <Navigate to="/auth" />;
  };

  return (
    <Router>
      <Routes>
        <Route
          path="/"
          element={requireAuth(<Main />)}
        />
        <Route
          path="/auth"
          element={isLoggedIn ?
                    <Navigate to="/" /> :
                    showLogin ?
                      <Login setShowLogin={setShowLogin} /> :
                      <Registration setShowLogin={setShowLogin} />}
        />
        <Route
          path="/session/:id"
          element={requireAuth(<Session />)}
        />
        <Route
          path="/session/:id/cards"
          element={requireAuth(<Cards />)}
        />
        <Route
          path="*"
          element={requireAuth(<NotFound />)}
        />
      </Routes>
    </Router>
  );
}

export default App;
