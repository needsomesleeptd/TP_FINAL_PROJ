import React, { useState, useEffect } from 'react';
import { useCookies } from 'react-cookie';
import './Main.css';


function Main() {
  const [cookies, removeCookie] = useCookies(['AccessToken']);
  const [participantsCount, setParticipantsCount] = useState('');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [focus, setFocus] = useState(0);
  const [sessionsData, setSessionsData] = useState([]);

  const handleFocus = (index) => {
    setFocus(index);
  };

  const userProfile = {
    username: 'Meet Match',
    avatar: 'https://cdn.icon-icons.com/icons2/38/PNG/512/search_4883.png'
  };

  const UserInfoRequest = async () => {
    try {
        const response = await fetch('http://localhost/api/sessions/getUser', {
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

  const createSession = async () => {
    try {
      const response = await fetch('http://localhost/api/sessions', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${cookies.AccessToken}`
          },
          body: JSON.stringify({
              "sessionName" : title,
              "sessionPeopleCap" : Number(participantsCount),
              "description" : description
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

  const handleTitleChange = (e) => {
    setTitle(e.target.value);
  };

  const handleDescriptionChange = (e) => {
    setDescription(e.target.value);
  };

  const handleParticipantsChange = (e) => {
    if (/^\d{0,3}$/.test(e.target.value)) {
      setParticipantsCount(e.target.value);
    }
  };

  const handleLogOut = () => {
    removeCookie("AccessToken");
  }

  const openModal = () => {
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
  };

  const joinSession = (sessionId) => {
    const sessionUrl = `http://localhost/api/sessions/${sessionId}`;
    window.location.href = sessionUrl;
  };
  
  const leaveSession = async (sessionId) => {
    try {
        const response = await fetch(`http://localhost/api/sessions/${sessionId}`, {
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

  const CreateSessionModal = ({ onCloseModal }) => {
    return (
      <div className="modal" style={{ display: 'block' }}>
        <div className="create-session-container modal-content">
          <span className="close" onClick={onCloseModal}>&times;</span>
          <h1>Создание встречи</h1>
          <div className="input-group">
            <input
              className="session-input"
              type="text"
              value={title}
              onChange={handleTitleChange}
              placeholder="Название"
              autoFocus={focus === 0}
              onFocus={() => handleFocus(0)}
            />
            <input
              className="session-input"
              type="text"
              value={description}
              onChange={handleDescriptionChange}
              placeholder="Описание"
              autoFocus={focus === 1}
              onFocus={() => handleFocus(1)}
            />
            <input
              className="session-input"
              type="text"
              value={participantsCount}
              onChange={handleParticipantsChange}
              placeholder="Количество участников"
              autoFocus={focus === 2}
              onFocus={() => handleFocus(2)}
            />
          </div>
          <button className="profile-button" onClick={createSession}>Создать</button>
        </div>
      </div>
    )
  }
  
  const ProfileHeader = ({ username, avatar, onOpenModal }) => {
    return (
      <div className="profile-header">
        <div className="profile-user">
          <img src={avatar} alt="User Avatar" className="profile-avatar" />
          <span className="profile-username">{username}</span>
        </div>
        {console.log("ff")}
        <div style={{display: "flex", gap: "10px"}}>
          <button className="profile-button" onClick={onOpenModal}>Создать встречу</button>
          <button className="profile-button" style={{marginRight: "50px"}} onClick={handleLogOut}>Выйти</button>
        </div>
      </div>
    );
  };
  
  const ProfileSession = ({ id, title, description, maxParticipants, participants, status }) => {
    return (
      <div className="profile-sessions session">
        <h3>{title}</h3>
        <p>{description}</p>
        <div>
          <span>{`Участники: ${participants}/${maxParticipants}`}</span>
          <span>{`Статус: ${status}`}</span>
        </div> 
        <button onClick={() => joinSession(id)}>Войти</button>
        <button onClick={() => leaveSession(id)}>Удалить</button>
      </div>
    );
  };

  return (
    <div className="app">
      {console.log("fff")}
      <ProfileHeader username={userProfile.username} avatar={userProfile.avatar} onOpenModal={openModal} />
      <div className="profile-content">
        <p className="profile-title">Ваши встречи</p>
        {sessionsData.length > 0 ? ( <div className="profile-sessions">
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
        </div> ) : (<p style={{marginLeft: "20px"}}>Нет встреч</p>)}
      </div>
      {isModalOpen && <CreateSessionModal onCloseModal={closeModal} />}
    </div>
  );
}

export default Main;
