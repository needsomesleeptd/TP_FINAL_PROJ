import React, { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import ConnectModal from './ConnectModal';
import InviteModal from './InviteModal';
import CreateModal from './CreateModal';
import { NavLink } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import Select from 'react-select';
import './Session.css'
import './Main.css'


const Session = (props) => {
  const { id } = useParams();
  const [sessionName, setSessionName] = useState('');
  const [sessionDesc, setSessionDesc] = useState(''); 
  const [maxParticipants, setMaxParticipants] = useState(0);
  const [participants, setParticipants] = useState([]);
  const [cookies, setCookie] = useCookies(['AccessToken', 'UserId', 'LoadedSession']);
  const [inputValue, setInputValue] = useState('');
  const [ready, setReady] = useState(false);
  const sessionId = id;
  const [showModal, setShowModal] = useState(false);
  const [showInviteModal, setShowInviteModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [date, setDate] = useState('');
  const navigate = useNavigate();

  const prevParticipants = useRef([]);
  const prevSessionName = useRef('');
  const prevSessionDesc = useRef('');
  const prevMaxParticipants = useRef(0);
  const prevDate = useRef('');

  const [selectedTags, setSelectedTags] = useState([]);
  const [showTagList, setShowTagList] = useState(false);
  const tagList = [
      { value: 'attractions', label: 'Достопримечательности' },
      { value: 'theater', label: 'Театры' },
      { value: 'restaurants', label: 'Рестораны' },
      { value: 'exhibition', label: 'Выставки' },
      { value: 'museums', label: 'Музеи' },
      { value: 'park', label: 'Парки' },
      { value: 'entertainment', label: 'Развлечения' },
      { value: 'education', label: 'Познавательно' },
      { value: 'tour', label: 'Туры' },
      { value: 'recreation', label: 'Отдых' },
      { value: 'questroom', label: 'Квеструмы' },
      { value: 'bar', label: 'Бары' },
      { value: 'cinema', label: 'Кинотеатры' },
      { value: 'concert', label: 'Концерты' },
      { value: 'clubs', label: 'Клубы' },
      { value: 'art-centers', label: 'Художественные центры' },
      { value: 'homesteads', label: 'Усадьбы' },
      { value: 'fountain', label: 'Фонтаны' },
    ]


  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const addTag = (tag) => {
    setSelectedTags([...selectedTags, tag]);
    setShowTagList(false);
  };
 
  const removeTag = (index) => {
    const updatedTags = [...selectedTags];
    updatedTags.splice(index, 1);
    setSelectedTags(updatedTags);
  };

  const link = window.location.href;

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

  const openEditModal = () => {
    setShowEditModal(true);
  };

  const closeEditModal = () => {
      setShowEditModal(false);
  };

  const handleEdit = (name, description, date, count) => {
    sessionModify(name, description, date, count);
    closeEditModal();
  };

  const handleSubmit = () => {
    patchSession(cookies.UserId);
  };

  const handleSubmit2 = () => {
    closeInviteModal();
  }

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

        console.log("polling");

        if (JSON.stringify(data.users) !== JSON.stringify(prevParticipants.current)) {
          setParticipants(data.users);
          prevParticipants.current = data.users;
        }

        if (data.sessionName !== prevSessionName.current) {
            setSessionName(data.sessionName);
            prevSessionName.current = data.sessionName;
        }

        if (data.description !== prevSessionDesc.current) {
            setSessionDesc(data.description);
            prevSessionDesc.current = data.description;
        }

        if (data.maxPeople !== prevMaxParticipants.current) {
            setMaxParticipants(data.maxPeople);
            prevMaxParticipants.current = data.maxPeople;
        }

        if (data.timeEnds) {
          const dd = data.timeEnds.split('T')[0];
          if (dd !== prevDate.current) {
            setDate(dd);
            prevDate.current = dd;
          }
        }

        const participant = data.users.find(participant => participant.ID === Number(cookies.UserId));
        
        if (participant && participant.Request !== '') {
          setInputValue(participant.Request);
          setSelectedTags(participant.Categories.map(value => tagList.find(item => item.value === value)));
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

      } catch (error) {
        console.error('Error creating session:', error);
      }
    };

    getSession();

    const pollingInterval = setInterval(getSession, 1000);

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

  const sessionModify = async (name, description, date, count) => {
    try {
      const response = await fetch('/api/sessions/update/'+ sessionId, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${cookies.AccessToken}`
          },
          body: JSON.stringify({
            'sessionID': sessionId,
            "sessionName": name,
            "sessionPeopleCap": count < participants.length ? participants.length : count,
            "description": description,
            "timeEnds": `${date}T23:59:00Z`
          })
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
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
            'newRequest': ready ? '' : inputValue.toString(),
            'newCategories': selectedTags.map((item) => item.value)
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

  const handleReadyClick = (event) => {
    event.preventDefault();
    putSession(cookies.meetmatchid);
    setReady(!ready);
  };

  const handleEditClick = (event) => {
    event.preventDefault();
    openEditModal();
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

  const handleCopyClick = (event) => {
    event.preventDefault();
    openInviteModal();
  }

  useEffect(() => {
    var objects = document.getElementsByClassName('precontainer');
    Array.from(objects).map((item) => {
      const img = new Image();
      img.src = item.getAttribute('data-src');
      img.onload = () => {
        sessionStorage.setItem("LoadedSession", true);
        item.style.backgroundImage = `url(${item.getAttribute('data-src')})`;
      };
      img.onerror = () => {
        sessionStorage.setItem("LoadedSession", false);
        console.error(`Error loading image: ${item.getAttribute('data-src')}`);
      };
    });

    objects = document.getElementsByClassName('logo');
    Array.from(objects).map((item) => {
      const img = new Image();
      img.src = item.getAttribute('data-src');
      img.onload = () => {
        item.src = item.getAttribute('data-src');
      };
      img.onerror = () => {
        console.error(`Error loading image: ${item.getAttribute('data-src')}`);
      };
    });
  }, []);

  const ComboBoxComponent = ({ options, onChange }) => {
    return (
      <Select
        className="custom-select"
        options={options}
        isSearchable
        placeholder={'Новый тег'}
        onChange={onChange}
        menuIsOpen={isMenuOpen}
        onMenuClose={() => setIsMenuOpen(false)}
        onMenuOpen={() => setIsMenuOpen(true)}
    />
    );
  };

  return (
    <div className={sessionStorage.getItem("LoadedSession") ? "precontainer loadedSession" : "precontainer"} data-src="/bg_session.png">
      <img src="data:image/gif;base64,R0lGODlhMgAbAIAAAP///wAAACH5BAEAAAEALAAAAAAyABsAAAIjjI+py+0Po5y02ouz3rz7D4biSJbmiabqyrbuC8fyTNf2zRUAOw==" data-src="/logo.png" class="logo" alt="Your Logo"></img>
      <ProfileHeader />
      <div class="container vertical-scroll-block2" style={{ height: "15vh", width: "95%", maxWidth: "800px" }}>
        <div class="container-info">
          <h2>{sessionName}</h2>
          {sessionDesc.length == '' ? null : <p>{sessionDesc}</p>}
        <form onSubmit={handleReadyClick} class="input-container">
          <input
            type="text"
            value={inputValue}
            onChange={handleInputChange}
            placeholder="Введите ваши пожелания..."
            disabled={ready}
            className="super-input"
            required
            />
          <div className="super-btns">
            <button class="profile-button" style={{width: "120px"}}>{ready ? "Не готов" : "Готов"}</button>
            <CopyToClipboard text={link}>
              <button onClick={(event) => handleCopyClick(event)} class="profile-button" style={{width: "120px"}} >Пригласить</button>
            </CopyToClipboard>
            <button onClick={(event) => handleEditClick(event)} class="profile-button" style={{width: "120px"}}>Изменить</button>
          </div>
        </form>
        <div className="input-container">
          {/* <p>Теги:</p> */}
          <div className="tags-container">
           {selectedTags.map((tag, index) => (
             <div key={index} className="tag">
               <p className="tags-p">{tag.label}</p>
               <button disabled={ready} className="tags-button" onClick={() => removeTag(index)}>✖</button>
             </div>
           ))}
          <div className="tag combobox-tag">
          {selectedTags.length !== tagList.length && (
            <ComboBoxComponent
            options={tagList.filter(tag => !selectedTags.some(selectedTag => selectedTag.value === tag.value))}
            onChange={addTag}
            />
            )}
          </div>
          </div>
          {/* {selectedTags.length !== tagList.length && (
          <div className="tags-new">
            <button onClick={() => setShowTagList(!showTagList)}>+</button>
            {showTagList && (
              <ul className="tag-list">
                {tagList.filter(tag => !selectedTags.includes(tag)).map((tag) => (
                  <li key={tag} onClick={() => addTag(tag)}>
                    {tag}
                  </li>
                ))}
              </ul>
            )}
          </div>
          )} */}
         </div>
        </div>
        {participants.length > 0 ? (
          <div>
            <p class="participants-count">Дата встречи: {date}</p>
            {
              participants.length < maxParticipants ?
              <p class="participants-count">Количество участников: {participants.length} / {maxParticipants}. Ждём пока остальные зайдут и будут готовы.</p>
              :
              <p class="participants-count">Количество участников: {participants.length} / {maxParticipants}. Ждём пока все будут готовы.</p>
            }
            <div>
            <table class="participants-table" style={{ width: "100%" }}>
              <thead>
                <tr>
                  <th>Пользователь</th>
                  <th>Готовность</th>
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
          </div>
        ) : (
          <p class="no-participants">Нет участников</p>
        )}
      </div>

      <ConnectModal showModal={showModal} sessionName={sessionName} sessionDesc={sessionDesc} sessionInfo={participants.length >= maxParticipants ? "Комната уже заполнена" : ""} handleUpload={handleSubmit} />
      <InviteModal showModal={showInviteModal} handleUpload={handleSubmit2} />
      { sessionName && (
        <CreateModal showModal={showEditModal} closeModal={closeEditModal} handleUpload={handleEdit}
        modalName="Изменение встречи" modalBtn="Сохранить"
        name_={sessionName} description_={sessionDesc} count_={maxParticipants} date_={date}
        />
      )}

      </div>
  );
};

export default Session;
