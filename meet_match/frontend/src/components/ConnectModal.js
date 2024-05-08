import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './CreateModal.css'

const ConnectModal = ({ showModal, sessionName, sessionDesc, handleUpload }) => {
  const navigate = useNavigate();  
  
  const joinSession = () => {
        handleUpload();
    }

    const cancelSession = () => {
      navigate('/');
    }

    if (!showModal) {
        return null;
    }

    return (
        <div className="upload-modal">
        <div className="upload-modal-content">
          <h1>Вы присоединяетесь к встрече "{sessionName}"</h1>
          <p>{sessionDesc}</p>
          <button className="modal-button" onClick={joinSession}>Продолжить</button>
          <p></p>
          <button className="modal-button" onClick={cancelSession}>На главную страницу</button>
        </div>
      </div>
    );
};

export default ConnectModal;
