import React, { useState, useEffect } from 'react';
import { useCookies } from 'react-cookie';
import { Link } from 'react-router-dom';
import CreateModal from './CreateModal';
import './Main.css';


const Main = () => {
  const [cookies, _, removeCookie] = useCookies(['AccessToken']);
  const [focus, setFocus] = useState(0);
  const [sessionsData, setSessionsData] = useState([]);
  const [showModal, setShowModal] = useState(false);

  const openModal = () => {
      setShowModal(true);
  };

  const closeModal = () => {
      setShowModal(false);
  };

  const handleFocus = (index) => {
    setFocus(index);
  };

  const userProfile = {
    username: 'Meet Match',
    avatar: 'https://w7.pngwing.com/pngs/665/132/png-transparent-user-defult-avatar-thumbnail.png'
  };

  const UserInfoRequest = async () => {
    try {
        const response = await fetch('/api/sessions/getUser', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin': '*',
                'Authorization': 'Bearer ' + cookies.AccessToken
            },
            body: JSON.stringify({
              'userID' : cookies.UserId
            })
        });

        if (!response.ok) {
          const errorMessage = await response.text();
          throw new Error(errorMessage);
        }

        const data = await response.json();
        setSessionsData(data.sessions ?? []);

    } catch (error) {
        alert(error.message);
    }
  }

  useEffect(() => {

    UserInfoRequest();

  }, [cookies]);

  const createSession = async (title, desc, count) => {
    try {
      const response = await fetch('/api/sessions', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${cookies.AccessToken}`
          },
          body: JSON.stringify({
              "sessionName" : title,
              "sessionPeopleCap" : Number(count),
              "description" : desc
          })
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const data = await response.json();
      const sessionId = data.sessionID;
      const sessionUrl = `http://localhost/session/${sessionId}`;
      window.location.href = sessionUrl;
    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  const handleLogOut = () => {
    removeCookie("AccessToken");
    removeCookie("UserId");
  }

  const joinSession = (sessionId) => {
    const sessionUrl = `http://localhost/sessions/${sessionId}`;
    window.location.href = sessionUrl;
  };
  
  const leaveSession = async (e, sessionId) => {
    e.stopPropagation();
    try {
        const response = await fetch(`/api/sessions/${sessionId}`, {
            method: 'DELETE',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${cookies.AccessToken}`
            },
            body: JSON.stringify({
                "sessionId" : sessionId
            })
        });
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        UserInfoRequest();
    } catch (error) {
      console.error('Error creating session:', error);
    }
  };
  
  const ProfileHeader = () => {
    return (
      <div className="profile-header">
        <Link to="/">Главная</Link>
        <Link to="/profile">Профиль</Link>
        <Link to="/about">О нас</Link>
      </div>
    );
  };

  const ProfileButtons = ({ onOpenModal }) => {
    return (
      <div className="profile-buttons">
        <button className="profile-button" onClick={onOpenModal}>Создать встречу</button>
        <button className="profile-button" onClick={handleLogOut}>Выйти</button>
      </div>
    );
  };
  
  const ProfileSession = ({ id, title, description, maxParticipants, participants, status }) => {
    return (
      <button className="profile-sessions session" onClick={() => joinSession(id)}>
        <div className="profile-posttitle">
          <h3>{title}</h3>
          <p>{description}</p>
        </div>
        <div>
          <p>{`Участники: ${participants}/${maxParticipants}`}</p>
          <p>{`Статус: ${status}`}</p>
          <p>{`Дата создания: 01.01.2000`}</p>
        </div>
        <button class="leave-button" onClick={(e) => leaveSession(e, id)}>
          X
        </button>
      </button>
    );
  };

  const VerticalScrollBlock = ({ children }) => {
    return (
      <div className="vertical-scroll-block">
        <div className="inner-scroll-content">{children}</div>
      </div>
    );
  }

  const handleUpload = (sessionName, sessionDesc, sessionCount) => {
      createSession(sessionName, sessionDesc, sessionCount);
      closeModal();
  };

  return (
    <div className="create-session-container-mega">
      <ProfileHeader />
      <ProfileButtons onOpenModal={openModal} />
      <div className="profile-content">
        {sessionsData.length > 0 ? (
          <div className="profile-sessions">
          <VerticalScrollBlock>
          {sessionsData.map((session, index) => (
            <ProfileSession
              key={index}
              id={session.sessionID}
              title={session.sessionName}
              description={session.description}
              maxParticipants={session.maxPeople}
              participants={session.users.length}
              status={session.status === 0 ? "Ожидание участников" :
                  (session.status == 1 ? "Просмотр карточек" : "Завершен")}
            />
          ))}
        </VerticalScrollBlock>
        </div>
        ) : (<p style={{marginLeft: "20px"}}>Нет встреч</p>)}
      </div>
      <CreateModal showModal={showModal} closeModal={closeModal} handleUpload={handleUpload} />
    </div>
  );
}

export default Main;
