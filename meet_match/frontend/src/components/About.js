import React, { useState, useEffect } from 'react';
import { NavLink } from 'react-router-dom';
import './NotFound.css'

function About() {
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
      <div style={{ width: "90%", maxWidth: "700px", textAlign: "center", padding: "20px", borderRadius: "10px", boxShadow: "0 2px 4px rgba(0, 0, 0, 0.1)" }}>
        <p><b>MeetMatch</b> - это инновационное приложение, которое делает выбор мест для встреч с друзьями увлекательным и простым. С помощью системы свайпов и машинного обучения, <b>MeetMatch</b> предлагает идеальные места, которые соответствуют предпочтениям всех участников группы.</p>
        <p><b>Как это работает?</b> Пользователь создает сессию и приглашает друзей. Каждый участник вводит свои предпочтения и выбирает категории мест, будь то кафе, парки, кинотеатры и т.д. Затем начинается свайпинг: пользователи оценивают предложенные варианты, и приложение анализирует их выборы с помощью модели машинного обучения.</p>
        <p>Основное преимущество <b>MeetMatch</b> - это учет предпочтений всех участников группы. Как только все участники оценят предложенные места, система находит те, которые понравятся всей компании, создавая идеальные условия для совместного досуга.</p>
        <p><b>MeetMatch</b> делает процесс планирования встреч не только эффективным, но и увлекательным, превращая его в интерактивную игру. Вместо долгих обсуждений и компромиссов, вы получаете быстрое и точное решение, которое учитывает вкусы каждого участника.</p>
        <p>С <b>MeetMatch</b> каждая встреча станет незабываемым событием, а процесс выбора мест - удовольствием. Попробуйте <b>MeetMatch</b> и убедитесь, что планирование досуга с друзьями может быть простым и веселым!</p>
      </div>
      <p></p>
      <p></p>
    </div>
  );
}

export default About;
