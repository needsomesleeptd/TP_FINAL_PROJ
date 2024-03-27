import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const CurSession = (props) => {
  const { id } = useParams();
  const [sessionId, setSessionId] = useState(id);
  const [participants, setParticipants] = useState([]);

  useEffect(() => {
    getSession();
    patchSession();
  }, [sessionId]);

  const getSession = async () => {
    try {
      const response = await fetch('http://localhost:8080/sessions/'+ sessionId, {
          method: 'POST',
          body: JSON.stringify({
            'sessionID': sessionId
          })
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const data = await response.json();
      setParticipants(data.UsersReqs);
      console.log(data.UsersReqs);
    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  const patchSession = async () => {
    try {
      const response = await fetch('http://localhost:8080/sessions/'+ sessionId, {
          method: 'FETCH',
          headers: {
            'Content-Type': 'application/json',
            // 'Origin': 'http://localhost:8080/sessions'
          },
          body: JSON.stringify({
            'sessionID': sessionId,
            'user': {
              'ID': 100,
              'Name': 'Sas',
              'Request': 'Koko'
            }
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

  const [inputValue, setInputValue] = useState('');

  const handleInputChange = (event) => {
    setInputValue(event.target.value);
  };

  const handleReadyClick = () => {
    alert(`Вы ввели: ${inputValue}`);
  };

  return (
    <div>
      <h2>Сессия</h2>
      <div>
        <input
          type="text"
          value={inputValue}
          onChange={handleInputChange}
          placeholder="Введите запрос..."
        />
        <button onClick={handleReadyClick}>Готов</button>
      </div>
      <p>Ссылка для приглашения: http://localhost:8080/sessions/{sessionId}</p>
      <p>Количество участников: {participants.length}</p>
      <table>
        <thead>
          <tr>
            <th>Пользователь</th>
            <th>Готов</th>
          </tr>
        </thead>
        <tbody>
          {participants.map((participant) => (
            <tr key={participant.id}>
              <td>{participant.Name}</td>
              <td>
                <input type="checkbox" disabled />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default CurSession;