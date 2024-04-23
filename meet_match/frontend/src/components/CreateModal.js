import React, { useState } from 'react';
import './CreateModal.css'

const CreateModal = ({ showModal, closeModal, handleUpload }) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [count, setCount] = useState(null);

    const createSession = () => {
        handleUpload(name, description, count);
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
        if (/^\d{0,3}$/.test(e.target.value)) {
          setCount(e.target.value);
        }
      };

    return (
        <div className="upload-modal">
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
            />
            <input
              className="session-input"
              type="text"
              value={description}
              onChange={(e) => handleDescriptionChange(e)}
              placeholder="Описание"
            />
            <input
              className="session-input"
              type="text"
              value={count}
              onChange={(e) => handleParticipantsChange(e)}
              placeholder="Количество участников"
            />
          </div>
          <button className="modal-button" onClick={createSession}>Создать</button>
        </div>
      </div>
    );
};

export default CreateModal;
