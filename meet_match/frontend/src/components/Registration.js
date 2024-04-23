import React, { useState, useEffect } from 'react';
import { RegisterRequest, LoginRequest } from './Auth';
import { useCookies } from 'react-cookie';
import './Auth.css'

const Registration = ({ setShowLogin }) => {
  const [name, setName] = useState('');
  const [age, setAge] = useState('');
  const [gender, setGender] = useState(true);
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [cookies, setCookie] = useCookies(['AccessToken']);
  
  const handleRegistration = async (event) => {
    event.preventDefault();
    const data = await RegisterRequest(name, Number(age), gender, login, password);
    if (data.status === "Error") {
      alert('Неверные данные');
    }
    else {
      const data2 = await LoginRequest(login, password);
      setCookie('AccessToken', data2.jwt, { path: '/' });
      setCookie('UserId', data2.userID, { path: '/' });
    } 
  };

  const handleButtonClick = (button, event) => {
    event.preventDefault();
    setGender(button);
  };

  const handleLoginRedirect = (event) => {
    event.preventDefault();
    setShowLogin(true);
  };

  return (
    <div class="container-auth">
      <div class="wrapper" style={{marginTop: "120px"}}>
        <link href='https://unpkg.com/boxicons@2.1.4/css/boxicons.min.css' rel='stylesheet'></link>
        <h1>Регистрация</h1>
        <form onSubmit={handleRegistration}>
          <div style={{display: "flex"}}>
            <div class="input-box">
              <input class="half-button" type="text" placeholder="Имя" value={name} onChange={(e) => setName(e.target.value)} required />
              <i class='bx bx-text'></i>
            </div>
            <div class="input-box">
              <input class="half-button" type="text" placeholder="Возраст" value={age} onChange={(e) => setAge(e.target.value)} required />
              <i class='bx bx-calendar' ></i>
            </div>
          </div>
          <div style={{display: "flex", margin: "0 7px 0 15px"}}>
            <p style={{width: "200px"}}>Выберите пол</p>
            <div className="rounded-rectangular-block">
              <button
                className={gender === true ? 'selected' : ''}
                onClick={(event) => handleButtonClick(true, event)}
              >
                <label>м</label>
              </button>
              <button
                className={gender === false ? 'selected' : ''}
                onClick={(event) => handleButtonClick(false, event)}
              >
                <label>ж</label>
              </button>
            </div>
          </div>
          <div class="input-box">
            <input type="text" placeholder="Логин" value={login} onChange={(e) => setLogin(e.target.value)} required />
            <i class='bx bxs-user'></i>
          </div>
          <div class="input-box">
            <input type="password" placeholder="Пароль" value={password} onChange={(e) => setPassword(e.target.value)} required />
            <i class='bx bxs-lock-alt' ></i>
          </div>
          <button type="submit" class="btn">Зарегистрироваться</button>
        </form>
        <div class="register-link">
          <p>Уже есть аккаунт? <a href="" onClick={handleLoginRedirect}>Авторизоваться</a></p>
        </div>
      </div>
    </div>
  );
};

export default Registration;
