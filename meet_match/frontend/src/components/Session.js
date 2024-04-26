import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import './Session.css'


const Session = (props) => {
  const { id } = useParams();
  const [sessionName, setSessionName] = useState([]);
  const [sessionDesc, setSessionDesc] = useState([]); 
  const [maxParticipants, setMaxParticipants] = useState(0);
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
        const response = await fetch('http://localhost/api/sessions/'+ sessionId, {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${cookies.AccessToken}`
            },
            body: JSON.stringify({
              'sessionID': sessionId
            })
        });
        const data = (await response.json()).session;
        setParticipants(data.users ?? []);
        setSessionName(data.sessionName);
        setSessionDesc(data.description);
        setMaxParticipants(data.maxPeople);
        const participant = data.users.find(participant => participant.ID === Number(cookies.UserId));
        if (participant.Request !== '') {
          setInputValue(participant.Request);
          setReady(true);
        }
        if (data.status === 1)
        {
          window.location.reload();
        }
      } catch (error) {
        console.error('Error creating session:', error);
      }
    };

    const pollingInterval = setInterval(getSession, 100);

    return () => clearInterval(pollingInterval);
  }, [cookies, sessionId]);

  const patchSession = async (id) => {
    try {
      const response = await fetch('http://localhost/api/sessions/'+ sessionId, {
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
    const participant = participants.find(participant => participant.ID === Number(cookies.UserId));
    try {
      const response = await fetch('http:/localhost/api/sessions/'+ sessionId, {
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
      <h1>Вы присоединяетесь к встрече "{sessionName}"</h1>
      <p>{sessionDesc}</p>
      <div class="input-container turbo-button">
        <button class="session-button turbo-button" onClick={handleSubmit}>Продолжить</button>
      </div>
      </div>
    );
  }

  return (
    <div class="container">
      <h2>{sessionName}</h2>
      <p>{sessionDesc}</p>
      <div class="input-container">
        <input
          type="text"
          value={inputValue}
          onChange={handleInputChange}
          placeholder="Введите ваши пожелания..."
          disabled={ready}
        />
        <button onClick={handleReadyClick}>{ready ? "Не готов" : "Готов"}</button>
      </div>
      <p class="invite-link">Ссылка для приглашения: http://redis-go-server:3000/session/{sessionId}</p>
      {participants.length > 0 ? (
        <div>
          <p class="participants-count">Количество участников: {participants.length} / {maxParticipants}</p>
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
