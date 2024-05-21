import React, { useState, useEffect } from 'react';
import { NavLink } from 'react-router-dom';
import './NotFound.css'

function Profile() {
  useEffect(() => {
    var objects = document.getElementsByClassName('error404');
    Array.from(objects).map((item) => {
      const img = new Image();
      img.src = item.getAttribute('data-src');
      img.onload = () => {
        sessionStorage.setItem("LoadedMain", true);
        item.style.backgroundImage = `url(${item.getAttribute('data-src')})`;
      };
      img.onerror = () => {
        sessionStorage.setItem("LoadedMain", false);
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
    <div className={sessionStorage.getItem("LoadedMain") ? "error404 loadedMain" : "error404"} data-src="/bg_main.png">
      <div>
        <img src="data:image/gif;base64,R0lGODlhMgAbAIAAAP///wAAACH5BAEAAAEALAAAAAAyABsAAAIjjI+py+0Po5y02ouz3rz7D4biSJbmiabqyrbuC8fyTNf2zRUAOw==" data-src="/logo.png" class="logo" alt="Your Logo"></img>
        <ProfileHeader />
      </div>
      <div style={{ width: "50%", textAlign: "center", padding: "20px", backgroundColor: "#24222240", borderRadius: "10px", boxShadow: "0 2px 4px rgba(0, 0, 0, 0.1)" }}>
        
      </div>
      <p></p>
      <p></p>
    </div>
  );
}

export default Profile;
