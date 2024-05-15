import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import { NavLink } from 'react-router-dom';
import './NotFound.css'

function NotFound() {
  const [cookies, setCookie] = useCookies(['AccessToken', 'LoadedMain']);
  const [countdown, setCountdown] = useState(3);
  const navigate = useNavigate();

  useEffect(() => {
    var objects = document.getElementsByClassName('error404');
    Array.from(objects).map((item) => {
      const img = new Image();
      img.src = item.getAttribute('data-src');
      img.onload = () => {
        setCookie("LoadedMain", true);
        item.style.backgroundImage = `url(${item.getAttribute('data-src')})`;
      };
      img.onerror = () => {
        setCookie("LoadedMain", false);
        console.error(`Error loading image: ${item.getAttribute('data-src')}`);
      };
    });

    objects = document.getElementsByClassName('logo');
    Array.from(objects).map((item) => {
      const img = new Image();
      img.src = item.getAttribute('data-src');
      img.onload = () => {
        item.src = item.getAttribute('data-src');
      };
      img.onerror = () => {
        console.error(`Error loading image: ${item.getAttribute('data-src')}`);
      };
    });
  }, []);

  const ProfileHeader = () => {
    return (
      <div className="profile-header">
        <NavLink to="/">Главная</NavLink>
        <NavLink to="/profile">Профиль</NavLink>
        <NavLink to="/about">О нас</NavLink>
      </div>
    );
  };

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
    <div className={cookies.LoadedMain ? "error404 loadedMain" : "error404"} data-src="/bg_main.png">
      <div>
        <img src="data:image/gif;base64,R0lGODlhMgAbAIAAAP///wAAACH5BAEAAAEALAAAAAAyABsAAAIjjI+py+0Po5y02ouz3rz7D4biSJbmiabqyrbuC8fyTNf2zRUAOw==" data-src="/logo.png" class="logo" alt="Your Logo"></img>
        <ProfileHeader />
      </div>
      <div className="fail404">
        <h1>Error 404</h1>
        <p>Страница не найдена</p>
      </div>
      <p></p>
    </div>
  );
}

export default NotFound;
