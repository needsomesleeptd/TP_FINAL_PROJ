import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import './Session.css'


const Session = (props) => {
  const { id } = useParams();
  const [participants, setParticipants] = useState([]);
  const [cookies] = useCookies(['AccessToken', 'UserId']);
  const [inputValue, setInputValue] = useState('');
  const [ready, setReady] = useState(false);
  const sessionId = id;

  const handleSubmit = () => {
    patchSession(cookies.UserId);
  };

  useEffect(() => {
    const getSession = async () => {
      try {
        const response = await fetch('http://localhost:8080/sessions/'+ sessionId, {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${cookies.AccessToken}`
            },
            body: JSON.stringify({
              'sessionID': sessionId
            })
        });
        const data = await response.json();
        setParticipants(data.UsersReqs ?? []);
        const participant = data.UsersReqs.find(participant => participant.ID == Number(cookies.UserId));
        if (participant.Request !== "") {
          setInputValue(participant.Request);
          setReady(true);
        }
        if (participants.length === 1 && data.UsersReqs.every(item => item.Request !== ''))
        {
          const sessionUrl = `http://localhost:3000/session/${sessionId}/cards`;
          window.location.href = sessionUrl; 
        }
      } catch (error) {
        console.error('Error creating session:', error);
      }
    };

    const pollingInterval = setInterval(getSession, 100);

    return () => clearInterval(pollingInterval);
  }, [cookies]);

  const patchSession = async (id) => {
    try {
      const response = await fetch('http://localhost:8080/sessions/'+ sessionId, {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${cookies.AccessToken}`
          },
          body: JSON.stringify({
            'sessionID': sessionId,
            'jwt': cookies.AccessToken
          })
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const data = await response.json();

      console.log(data);
    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  const putSession = async (id) => {
    const participant = participants.find(participant => participant.ID == Number(cookies.UserId));
    try {
      const response = await fetch('http://localhost:8080/sessions/'+ sessionId, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${cookies.AccessToken}`
          },
          body: JSON.stringify({
            'sessionID': sessionId,
            'userIDToModify': Number(cookies.UserId),
            'newName': participant.Name,
            'newRequest': ready ? '' : inputValue.toString()
          })
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const data = await response.json();

      console.log(data);
    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  const handleInputChange = (event) => {
    setInputValue(event.target.value);
  };

  const handleReadyClick = () => {
    putSession(cookies.meetmatchid);
    setReady(!ready);
  };

  console.log()
  if (participants.length > 0 && !participants.find(participant => participant.ID === Number(cookies.UserId))) {
    return (
      <div className="create-session-container">
      <h1>Вход в сессию</h1>
      <button class="session-button" onClick={handleSubmit}>Продолжить</button>
      </div>
    );
  }

  return (
    <div class="container">
      <h2>Сессия</h2>
      <div class="input-container">
        <input
          type="text"
          value={inputValue}
          onChange={handleInputChange}
          placeholder="Введите запрос..."
          disabled={ready}
        />
        <button onClick={handleReadyClick}>{ready ? "Не готов" : "Готов"}</button>
      </div>
      <p class="invite-link">Ссылка для приглашения: http://localhost:8080/sessions/{sessionId}</p>
      {participants.length > 0 ? (
        <div>
          <p class="participants-count">Количество участников: {participants.length} / 2</p>
          <table class="participants-table">
            <thead>
              <tr>
                <th>Пользователь</th>
                <th>Готов</th>
              </tr>
            </thead>
            <tbody>
              {participants.map((participant) => (
                <tr key={participant.ID}>
                  <td>{participant.Name}</td>
                  <td>
                    <label class="checkbox-container">
                      <input type="checkbox" class="checkbox-input" disabled checked={participant.Request !== ''} />
                      <span class="checkbox-custom"></span>
                    </label>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <p class="no-participants">Нет участников</p>
      )}
    </div>
  );
};

export default Session;
