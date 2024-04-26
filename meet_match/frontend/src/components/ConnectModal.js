import React, { useState } from 'react';
import './CreateModal.css'

const ConnectModal = ({ showModal, sessionName, sessionDesc, handleUpload }) => {
    const joinSession = () => {
        handleUpload();
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
        </div>
      </div>
    );
};

export default ConnectModal;
