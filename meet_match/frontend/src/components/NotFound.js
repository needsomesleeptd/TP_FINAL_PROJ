import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './NotFound.css'

function NotFound() {
  const [countdown, setCountdown] = useState(3);
  const navigate = useNavigate();

  // useEffect(() => {
  //   const timer = setInterval(() => {
  //     setCountdown((prevCount) => prevCount - 1);
  //   }, 1000);

  //   setTimeout(() => {
  //     clearInterval(timer);
  //     navigate('/');
  //   }, 3000);

  //   return () => {
  //     clearInterval(timer);
  //   };
  // }, []);

  return (
    <div className="error404">
      <h1>Error 404</h1>
      <p>Страница не найдена</p>
    </div>
  );
}

export default NotFound;
