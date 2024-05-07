import React, { useState } from 'react';
import './CreateModal.css'

const InviteModal = ({ showModal, handleUpload }) => {
    const okSession = () => {
        handleUpload();
    }

    if (!showModal) {
        return null;
    }

    return (
        <div className="upload-modal">
        <div className="upload-modal-content">
          <h1>Вы скопировали ссылку-приглашение!</h1>
          <p>Отправьте её своим друзьям, которых хотите позвать на встречу.</p>
          <button className="modal-button" onClick={okSession}>Ок</button>
        </div>
      </div>
    );
};

export default InviteModal;
