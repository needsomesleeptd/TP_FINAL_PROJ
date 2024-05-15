import React, { useState, useEffect } from 'react';
import { RegisterRequest, LoginRequest } from './Auth';
import { useCookies } from 'react-cookie';
import './Auth.css'

const Registration = ({ setShowLogin }) => {
  const [name, setName] = useState('');
  const [age, setAge] = useState();
  const [gender, setGender] = useState(true);
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [cookies, setCookie] = useCookies(['AccessToken', 'LoadedMain']);

  useEffect(() => {
    var objects = document.getElementsByClassName('container-auth');
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

  const handleAgeChange = (e) => {
    const num = parseInt(e.target.value, 10);
    if (num < 1) {
      setAge(1);
    } else if (num > 100) {
      setAge(100);
    } else {
      setAge(num);
    }
  };

  return (
    <div className={cookies.LoadedMain ? "container-auth loadedMain" : "container-auth"} data-src="/bg_main.png">
      <img src="data:image/gif;base64,R0lGODlhMgAbAIAAAP///wAAACH5BAEAAAEALAAAAAAyABsAAAIjjI+py+0Po5y02ouz3rz7D4biSJbmiabqyrbuC8fyTNf2zRUAOw==" data-src="/logo.png" class="logo" alt="Your Logo"></img>
      <div class="wrapper">
        <link href='https://unpkg.com/boxicons@2.1.4/css/boxicons.min.css' rel='stylesheet'></link>
        <h1>Регистрация</h1>
        <form onSubmit={handleRegistration}>
          <div style={{display: "flex"}}>
            <div class="input-box">
              <input class="half-button" type="text" placeholder="Имя" value={name} onChange={(e) => setName(e.target.value)} required />
              <i class='bx bx-text'></i>
            </div>
            <div class="input-box">
              <input class="half-button" type="number" placeholder="Возраст" value={age} min={1} max={100} onChange={(e) => handleAgeChange(e)} required />
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
      <p></p>
    </div>
  );
};

export default Registration;
