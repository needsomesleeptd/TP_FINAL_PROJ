import React, { useState } from 'react';


function Session() {

  const createSession = async () => {
    try {
      const response = await fetch('http://localhost:8080/sessions', {
          method: 'POST'
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const data = await response.json();
      const sessionId = data.sessionID;
      const sessionUrl = `http://localhost:3000/session/${sessionId}`;
      window.location.href = sessionUrl;
    } catch (error) {
      console.error('Error creating session:', error);
    }
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
