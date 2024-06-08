import React, { useEffect, useState, useContext  } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import { useParams } from 'react-router-dom';
import Registration from './components/Registration'
import Login from './components/Login'
import Main from './components/Main'
import Session from './components/Session'
import Cards from './components/Cards'
import Match from './components/Match'
import About from './components/About';
import Profile from './components/Profile';
import NotFound from './components/NotFound'
import './custom.css'
import { UserProvider, UserContext } from './components/MyContext';

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
    console.log(id);
    const [status, setStatus] = useState('');
    const {letsSwipe} = useContext(UserContext);

    useEffect(() => {
      const CheckSession = async (sessionId) => {
        try {
          const response = await fetch('/api/sessions/'+ id, {
              method: 'POST',
              headers: {
                'Authorization': `Bearer ${cookies.AccessToken}`
              },
              body: JSON.stringify({
                'sessionID': id
              })
          });
          const data = await response.json();
          console.log(data);
          if (data.Response.status === "OK") {
            setStatus(data.session.status);
            if (!data.session.users.map(u => u.ID).find(m => m === cookies.UserId)) {
              setStatus(0);
            }
          }
          else {
            setStatus(-1);
          }
        } catch (error) {
          setStatus(-1);
          console.error('Error creating session:', error);
        }
      };
  
      CheckSession();

      const interval = setInterval(() => {
        CheckSession(); // Fetch the status at regular intervals
      }, 1000); // Fetch every second
  
      return () => clearInterval(interval); // Clean up the interval on component unmount
    }, [id]);

    console.log(`status: ${status}`);
  
    if (status === 0) {
      return <Session />;
    } else if (status === 1) {
      return <Cards />;
    } else if (status >= 2) {
      return letsSwipe ? <Cards /> : <Match />;
    } else if (status === -1) {
      return <NotFound />;
    } else {
      return <div></div>;
    }
  };


  return (
    <UserProvider>
    <Router>
      <Routes>
        <Route
          path="/"
          element={requireAuth(<Main />)}
        />
        <Route
          path="/about"
          element={requireAuth(<About />)}
        />
        <Route
          path="/profile"
          element={requireAuth(<Profile />)}
        />
        <Route
          path="/auth"
          element={(isLoggedIn && hasUserId) ?
            <Navigate to='/' /> :
                    showLogin ?
                    <Login setShowLogin={setShowLogin} /> :
                    <Registration setShowLogin={setShowLogin} />}
                    />
        <Route
          path="/session/:id"
          element={(isLoggedIn && hasUserId) ?
            <DataFetcher /> :
            showLogin ?
            <Login setShowLogin={setShowLogin} /> :
            <Registration setShowLogin={setShowLogin} />}
            />
        <Route
          path="*"
          element={requireAuth(<NotFound />)}
          />
      </Routes>
    </Router>
    </UserProvider>
  );
}

export default App;
