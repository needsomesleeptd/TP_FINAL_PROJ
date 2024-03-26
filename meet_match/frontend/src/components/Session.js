import React, { useState } from 'react';


function Session() {

  const createSession = async () => {
    // const response = await fetch('https://localhost:3000/sessions', {
    //     method: 'PUSH',
    //     headers: {
    //         'Content-Type': 'application/json'
    //     }
    // });
    // if (response.ok){
        // const data = await response.json();
        // const sessionId = data.id;
        const sessionId = 1;
        const sessionUrl = `http://localhost:3000/session/${sessionId}`;
        window.location.href = sessionUrl;
    // }
  };
  
  const joinSession = async () => {
    // const response = await fetch('https://localhost:3000/sessions', {
    //     method: 'PUSH',
    //     headers: {
    //         'Content-Type': 'application/json'
    //     }
    // });
    // if (response.ok){
        // const data = await response.json();
        // const sessionId = data.id;
        const sessionId = 1;
        const sessionUrl = `http://localhost:3000/session/${sessionCode}`;
        window.location.href = sessionUrl;
    // }
  };

  const [sessionCode, setSessionCode] = useState('');

  const handleInputChange = (event) => {
    setSessionCode(event.target.value);
  };

  return (
    <div>
        <h1>Главная страница</h1>
        <div>
            <button onClick={createSession}>Создать сессию</button>
        </div>
        <p></p>
        <div>
        <input
            type="text"
            placeholder="Введите код сессии"
            value={sessionCode}
            onChange={handleInputChange}/> 
            <button onClick={joinSession}>Подключиться</button>
        </div>
    </div>
  );
}

export default Session;
