import React, { useState, useEffect, lazy, Suspense } from 'react';
import { LoginRequest } from './Auth';
import { useCookies } from 'react-cookie';
import './Auth.css'


const Login = ({ setShowLogin }) => {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [cookies, setCookie] = useCookies(['AccessToken', 'LoadedMain']);
  const [errorMessage, setErrorMessage] = useState('');
  
  useEffect(() => {
    var objects = document.getElementsByClassName('container-auth');
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

  const handleLogin = async (event) => {
    event.preventDefault();
    try {
      const data = await LoginRequest(login, password);
      if (data.status === "Error") {
        setErrorMessage('Неправильно указан логин и/или пароль');
      }
      else {
        setCookie('AccessToken', data.jwt, { path: '/' });
        setCookie('UserId', data.userID, { path: '/' });
      }
    } catch (error) {
      setErrorMessage('Непредвиденная ошибка. Попробуйте зайти попозже.');
    }
  };

  const handleRegistrationRedirect = (event) => {
    event.preventDefault();
    setShowLogin(false);
  };

  return (
    <div className={sessionStorage.getItem("LoadedMain") ? "container-auth loadedMain" : "container-auth"} data-src="/bg_main.png">
      <img src="data:image/gif;base64,R0lGODlhMgAbAIAAAP///wAAACH5BAEAAAEALAAAAAAyABsAAAIjjI+py+0Po5y02ouz3rz7D4biSJbmiabqyrbuC8fyTNf2zRUAOw==" data-src="/logo.png" class="logo" alt="Your Logo"></img>
      <div class="wrapper">
        <link href='https://unpkg.com/boxicons@2.1.4/css/boxicons.min.css' rel='stylesheet'></link>
        <h1>Авторизация</h1>
        <form onSubmit={handleLogin}>
          <div class="input-box">
            <input type="text" placeholder="Логин" value={login} onChange={(e) => setLogin(e.target.value)} required />
            <i class='bx bxs-user'></i>
          </div>
          <div class="input-box">
            <input type="password" placeholder="Пароль" value={password} onChange={(e) => setPassword(e.target.value)} required />
            <i class='bx bxs-lock-alt' ></i>
          </div>
          {/* <div class="remember-forgot">
            <label></label>
            <a href="#">Забыли пароль?</a>
          </div> */}
          {errorMessage && (
          <p class="error-message" style={{ color: 'red' }}>{errorMessage}</p>
          )}
          <button type="submit" class="btn">Войти</button>
        </form>
        <div class="register-link">
          <p>Нет аккаунта? <a href="" onClick={handleRegistrationRedirect}>Зарегистрироваться</a></p>
        </div>
      </div>
      <p></p>
    </div>
  );
};

export default Login;
