import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import { NavLink } from 'react-router-dom';
import ConnectModal from './ConnectModal';
import './Session.css'
import './Main.css'
import InviteModal from './InviteModal';

const Session = (props) => {
  const { id } = useParams();
  const [sessionName, setSessionName] = useState([]);
  const [sessionDesc, setSessionDesc] = useState([]); 
  const [maxParticipants, setMaxParticipants] = useState(0);
  const [participants, setParticipants] = useState([]);
  const [cookies] = useCookies(['AccessToken', 'UserId']);
  const [inputValue, setInputValue] = useState('');
  const [ready, setReady] = useState(false);
  const sessionId = id;
  const [showModal, setShowModal] = useState(false);
  const [showInviteModal, setShowInviteModal] = useState(false);

  const openModal = () => {
      setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
  };

  const openInviteModal = () => {
    setShowInviteModal(true);
  };

  const closeInviteModal = () => {
    setShowInviteModal(false);
  };

  const handleSubmit = () => {
    patchSession(cookies.UserId);
  };

  const handleSubmit2 = () => {
    closeInviteModal();
  };


  useEffect(() => {
    const getSession = async () => {
      try {
        const response = await fetch('/api/sessions/'+ sessionId, {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${cookies.AccessToken}`
            },
            body: JSON.stringify({
              'sessionID': sessionId
            })
        });
        const data = (await response.json()).session;
        setParticipants(data.users ?? []);
        setSessionName(data.sessionName);
        setSessionDesc(data.description);
        setMaxParticipants(data.maxPeople);
        const participant = data.users.find(participant => participant.ID === Number(cookies.UserId));
        console.log(participant);
        
        if (participant && participant.Request !== '') {
          setInputValue(participant.Request);
          setReady(true);
        }


        if (data.users.length > 0 && !participant)
        {
          openModal();
        }
        else
        {
          closeModal();
        }

        if (data.status === 1)
        {
          window.location.reload();
        }
      } catch (error) {
        console.error('Error creating session:', error);
      }
    };

    getSession();

    const pollingInterval = setInterval(getSession, 500);

    return () => clearInterval(pollingInterval);
  }, [cookies, sessionId]);

  const patchSession = async (id) => {
    try {
      const response = await fetch('/api/sessions/'+ sessionId, {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${cookies.AccessToken}`
          },
          body: JSON.stringify({
            'sessionID': sessionId,
            'jwt': cookies.AccessToken
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

  const putSession = async (id) => {
    const participant = participants.find(participant => participant.ID === Number(cookies.UserId));
    try {
      const response = await fetch('/api/sessions/'+ sessionId, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${cookies.AccessToken}`
          },
          body: JSON.stringify({
            'sessionID': sessionId,
            'userIDToModify': Number(cookies.UserId),
            'newName': participant.Name,
            'newRequest': ready ? '' : inputValue.toString()
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

  const handleInputChange = (event) => {
    setInputValue(event.target.value);
  };

  const handleReadyClick = () => {
	if (!inputValue.toString()) 
	{
   	return;
	}
	putSession(cookies.meetmatchid);
    setReady(!ready);
  };

  const ProfileHeader = () => {
    return (
      <div className="profile-header">
        <NavLink to="/">Главная</NavLink>
        <NavLink to="/profile">Профиль</NavLink>
        <NavLink to="/about">О нас</NavLink>
      </div>
    );
  };

  const handleCopyClick = () => {
    navigator.clipboard.writeText(window.location.href);
    openInviteModal();
  }

  return (
    <div class="precontainer">
      <ProfileHeader />
      <div class="container">
        <div class="container-info">
          <h2>{sessionName}</h2>
          <p>{sessionDesc}</p>
        <div class="input-container">
          <input
            type="text"
            value={inputValue}
            onChange={handleInputChange}
            placeholder="Введите ваши пожелания..."
            disabled={ready}
            />
          <button onClick={handleReadyClick} class="profile-button" style={{width: 150}}>{ready ? "Не готов" : "Готов"}</button>
  
        </div>
        <p class="invite-link">Отправьте ссылку друзьям: {window.location.href}. Ждём других участников. Все должны ввести пожелания и нажать на "Готов".</p>
        </div>
        {participants.length > 0 ? (
          <div>
            <p class="participants-count">Количество участников: {participants.length} / {maxParticipants}</p>
            <table class="participants-table">
              <thead>
                <tr>
                  <th>Пользователь</th>
                  <th>Готов</th>
                </tr>
              </thead>
              <tbody>
                {participants.map((participant) => (
                  <tr key={participant.ID}>
                    <td>{participant.Name}</td>
                    <td>
                      <label class="checkbox-container">
                        <input type="checkbox" class="checkbox-input" disabled checked={participant.Request !== ''} />
                        <span class="checkbox-custom"></span>
                      </label>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p class="no-participants">Нет участников</p>
        )}
      </div>

      <ConnectModal showModal={showModal} sessionName={sessionName} sessionDesc={sessionDesc} handleUpload={handleSubmit} />
      <InviteModal showModal={showInviteModal} handleUpload={handleSubmit2} />

      </div>
  );

};

export default Session;
