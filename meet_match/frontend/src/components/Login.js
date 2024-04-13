import React, { useState, useEffect } from 'react';
import { LoginRequest } from './Auth';
import { useCookies } from 'react-cookie';
import './Auth.css'


const Login = ({ setShowLogin }) => {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [cookies, setCookie] = useCookies(['AccessToken']);

  const handleLogin = async (event) => {
    event.preventDefault();
    const data = await LoginRequest(login, password);
    setCookie('AccessToken', data.jwt, { path: '/' });
    setCookie('UserId', data.userID, { path: '/' });
    window.location.href = '/';
  };

  const handleRegistrationRedirect = (event) => {
    event.preventDefault();
    setShowLogin(false);
  };

  return (
    <div class="container-auth">
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
          <div class="remember-forgot">
            <label><input type="checkbox" /> Запомнить меня</label>
            <a href="#">Забыли пароль?</a>
          </div>
          <button type="submit" class="btn">Войти</button>
        </form>
        <div class="register-link">
          <p>Нет аккаунта? <a href="" onClick={handleRegistrationRedirect}>Зарегистрироваться</a></p>
        </div>
      </div>
    </div>
  );
};

export default Login;
