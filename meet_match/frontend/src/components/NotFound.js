import React, { useState, useEffect } from 'react';

function NotFound() {
  const [countdown, setCountdown] = useState(3);

  useEffect(() => {
    const timer = setInterval(() => {
      setCountdown((prevCount) => prevCount - 1);
    }, 1000);

    setTimeout(() => {
      clearInterval(timer);
      window.location.href = '/';
    }, 3000);

    return () => {
      clearInterval(timer);
    };
  }, []);

  return (
    <div>
      <h1>Error 404</h1>
      <p>Страница не найдена</p>
      <p>Переход на главную страницу через {countdown} сек.</p>
    </div>
  );
}

export default NotFound;
