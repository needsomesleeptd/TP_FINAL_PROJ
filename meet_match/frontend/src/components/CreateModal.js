import React, { useState } from 'react';
import './CreateModal.css'

const CreateModal = ({ showModal, closeModal, handleUpload }) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [count, setCount] = useState(2);

    const createSession = (event) => {
        event.preventDefault();
        if (Number(count) <= 0)
        {
          handleUpload(name, description, 1);
        }
        else
        {
          handleUpload(name, description, count);
        }
    }

    if (!showModal) {
        return null;
    }

    const handleTitleChange = (e) => {
        setName(e.target.value);
      };
    
      const handleDescriptionChange = (e) => {
        setDescription(e.target.value);
      };
    
      const handleParticipantsChange = (e) => {
        const num = parseInt(e.target.value, 10);
        if (num < 1) {
          setCount(1);
        } else if (num > 100) {
          setCount(100);
        } else {
          setCount(num);
        }
      };

    return (
        <div className="upload-modal">
        <form onSubmit={createSession}>
        <div className="upload-modal-content">
          <span className="close" onClick={closeModal}>&times;</span>
          <h1>Создание встречи</h1>
          <div className="input-group">
            <input
              className="session-input"
              type="text"
              value={name}
              onChange={(e) => handleTitleChange(e)}
              placeholder="Название"
              required
            />
            <input
              className="session-input"
              type="text"
              value={description}
              onChange={(e) => handleDescriptionChange(e)}
              placeholder="Сообщение для участников"
            />
            <div className="input-group" style={{ display: 'flex', alignItems: 'center', width: '100%', justifyContent: "space-between" }}>
            <p>Количество участников:</p>
              <input
                  type="number"
                  className="session-input"
                  value={count}
                  onChange={(e) => handleParticipantsChange(e)}
                  min={1}
                  max={100}
                  style={{ width: '25%' }}
                  required
              />
          </div>
          <input
                  type="range"
                  id="participantRange"
                  name="participantRange"
                  min={1}
                  max={100}
                  value={count}
                  onChange={(e) => handleParticipantsChange(e)}
                  style={{ width: '100%' }}
              />
          </div>
          <button className="modal-button">Создать</button>
        </div>
        </form>
      </div>
    );
};

export default CreateModal;
