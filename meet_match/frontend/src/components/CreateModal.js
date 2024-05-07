import React, { useState } from 'react';
import './CreateModal.css'

const CreateModal = ({ showModal, closeModal, handleUpload }) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [count, setCount] = useState(null);

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
        if (/^\d*\.?\d*$/.test(e.target.value)) {
          setCount(e.target.value);
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
            <input
              className="session-input"
              type="text"
              value={count}
              onChange={(e) => handleParticipantsChange(e)}
              placeholder="Количество участников"
		required
            />
          </div>
          <button className="modal-button">Создать</button>
        </div>
	</form>
      </div>
    );
};

export default CreateModal;
