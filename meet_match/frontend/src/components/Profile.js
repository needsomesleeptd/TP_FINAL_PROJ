import React, { useState, useEffect } from 'react';
import { NavLink } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import { motion } from 'framer-motion';
import './NotFound.css'

function Profile() {
  const [cookies] = useCookies(['AccessToken']);
  const [info, setInfo] = useState(null);
  const [cards, setCards] = useState([]);
  const [selectedCard, setSelectedCard] = useState(-1);

  useEffect(() => {
    var objects = document.getElementsByClassName('error404');
    Array.from(objects).map((item) => {
      const img = new Image();
      img.src = item.getAttribute('data-src');
      img.onload = () => {
        sessionStorage.setItem("LoadedMain", true);
        item.style.backgroundImage = `url(${item.getAttribute('data-src')})`;
      };
      img.onerror = () => {
        sessionStorage.setItem("LoadedMain", false);
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

  const ProfileHeader = () => {
    return (
      <div className="profile-header">
        <NavLink to="/">Главная</NavLink>
        <NavLink to="/profile">Профиль</NavLink>
        <NavLink to="/about">О нас</NavLink>
      </div>
    );
  };

  useEffect(() => {
    const getStats = async () => {
      try {
        const response = await fetch(`/api/user/stats`, {
            method: 'GET',
            headers: {
              'Authorization': `Bearer ${cookies.AccessToken}`
            }
        });
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
  
        const data = await response.json();
        console.log(data);

        setInfo(data.peron_stats);
        setCards([data.peron_stats.most_liked_place, data.peron_stats.most_disliked_place]);
    
      } catch (error) {
        console.error('Error creating session:', error);
      }
    };
  
    getStats();

   }, []);

   const handleCardClick = (index) => {
    console.log(index);
    setSelectedCard(selectedCard === index ? -1 : index);
  };

  return (
    <div className={sessionStorage.getItem("LoadedMain") ? "error404 loadedMain" : "error404"} data-src="/bg_main.png">
      <div>
        <img src="data:image/gif;base64,R0lGODlhMgAbAIAAAzP///wAAACH5BAEAAAEALAAAAAAyABsAAAIjjI+py+0Po5y02ouz3rz7D4biSJbmiabqyrbuC8fyTNf2zRUAOw==" data-src="/logo.png" class="logo" alt="Your Logo"></img>
        <ProfileHeader />
      </div>
      {info && <div style={{ width: "50%", textAlign: "center", padding: "20px", borderRadius: "10px", boxShadow: "0 2px 4px rgba(0, 0, 0, 0.1)", maxWidth: "500px", display: "flex", flexDirection: "column", alignItems: "center" }}>
        <p>Количество текущих сессий: {info.sessions_count}</p>
        <p>Всего свайпов: {info.personal_stats.swiped} ({info.personal_stats.positive_swipes} лайков / {info.personal_stats.negative_swipes} дизлайков)</p>
        <div style={{display: "flex", flexDirection: "row", justifyContent: "space-between" }}>
        {cards.slice().reverse().map((card, index) => (
          <div>
          <p style={{
            marginLeft: index === 1 ? "-20px" : "auto",
            marginRight: index === 1 ? "auto" : "-20px"
          }}>{index === 0 ? `Самая любимая карта (${info.most_liked_scrolled_count}):` : `Самая нелюбимая карта (${info.most_disliked_scrolled_count}):`}</p>
          <motion.div
          key={index}
          onClick={() => handleCardClick(index)}
          animate={{ scale: selectedCard === index ? 0.85 : 0.8 }}
          class="cards-card"
          transition={{ duration: 0.3, type: 'spring', stiffness: 300 }}
          style={{ position: "relative", marginTop: "-40px",
          marginLeft: index === 1 ? "-20px" : 0,
          marginRight: index === 1 ? 0 : "-20px"
        }}
          >
          <div style={{display: "flex", flexDirection: "column", justifyContent: "space-between" }}>
            {selectedCard === index ? (
              <>
                <p style={{textAlign: "center", margin: "20px 10px", fontSize: "16px", fontWeight: "bold" }}>{card.title}</p>
                <p className="cards-descr" dangerouslySetInnerHTML={{ __html: card.description }} />
                {card.age_restriction && <p style={{marginLeft: "10px", textAlign: "left", marginBottom: "10px", fontSize: "14px" }}>Возраст: {card.age_restriction}</p>}
                {card.cost && <p style={{marginLeft: "10px", textAlign: "left", marginBottom: "10px", fontSize: "14px" }}>Цена: {card.cost}</p>}
                {card.timetable && <p style={{marginLeft: "10px", textAlign: "left", marginBottom: "10px", fontSize: "14px" }}>Расписание: {card.timetable}</p>}
                {card.subway && <p style={{marginLeft: "10px", textAlign: "left", marginBottom: "10px", fontSize: "14px" }}>Метро: {card.subway}</p>}
                {card.site_url && <p style={{marginLeft: "10px",  textAlign: "left", marginBottom: "10px", fontSize: "12px" }}>Сайт: <a href={card.site_url} target="_blank" rel="noopener noreferrer"
                onClick={(e) => {
                  e.stopPropagation();
                }}
                style={{color: "white"}}
                >*Кликай*</a></p>}
              </>
            ) :
            <>
            <img src={card.image} alt="" class="cards-img" />
            <p style={{textAlign: "center", margin: "10px", fontSize: "16px" }}>{card.title}</p>
            <p style={{textAlign: "center", margin: "10px", fontSize: "12px" }}>*Нажмите, чтобы узнать подробнее*</p>
            </>
            }
          </div>
        </motion.div>
        </div>
        ))}
      </div>
      </div> }
      <p></p>
      <p></p> 
    </div>
  );
}

export default Profile;
