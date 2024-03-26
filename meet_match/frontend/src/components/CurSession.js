import React, { useState } from 'react';

const CurSession = () => {
  const [participants, setParticipants] = useState([
    { id: 1, name: 'Вася', checked: false },
    { id: 2, name: 'Петя', checked: false },
    { id: 3, name: 'Гриша', checked: false },
  ]);

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
      <p>Ссылка для приглашения: localhost:3000/session/fJy1F3j</p>
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
              <td>{participant.name}</td>
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