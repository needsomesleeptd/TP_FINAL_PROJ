import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import { useParams } from 'react-router-dom';
import Registration from './components/Registration'
import Login from './components/Login'
import Main from './components/Main'
import Session from './components/Session'
import Cards from './components/Cards'
import Match from './components/Match'
import NotFound from './components/NotFound'
import './custom.css'

function App() {
  const [cookies] = useCookies(['AccessToken', 'UserId']);
  const isLoggedIn = !!cookies.AccessToken;
  const hasUserId = !!cookies.UserId;
  const [showLogin, setShowLogin] = useState(true);

  const requireAuth = (element) => {
    return (isLoggedIn && hasUserId) ? element : <Navigate to="/auth" />;
  };

  const DataFetcher = ({ sessionId }) => {
    const { id } = useParams();
    const [status, setStatus] = useState('');

    useEffect(() => {
      const CheckSession = async (sessionId) => {
        try {
          const response = await fetch('http://localhost:8080/sessions/'+ id, {
              method: 'POST',
              headers: {
                'Authorization': `Bearer ${cookies.AccessToken}`
              },
              body: JSON.stringify({
                'sessionID': id
              })
          });
          const data = await response.json();
          if (data.Response.status === "OK") {
            setStatus(data.session.status);
          }
          else {
            setStatus(-1);
          }
        } catch (error) {
          console.error('Error creating session:', error);
        }
      };
  
      CheckSession();
    }, [sessionId]);
  
    if (status === 0) {
      return <Session />;
    } else if (status === 1) {
      return <Cards />;
    } else if (status === 2) {
      return <Match />;
    } else if (status === -1) {
      return <NotFound />;
    } else {
      return <div></div>;
    }
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
          element={(isLoggedIn && hasUserId) ?
                    <Navigate to="/" /> :
                    showLogin ?
                      <Login setShowLogin={setShowLogin} /> :
                      <Registration setShowLogin={setShowLogin} />}
        />
        <Route
          path="/session/:id"
          element={<DataFetcher />}
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
