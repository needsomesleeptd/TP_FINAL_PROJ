import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useParams } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import { NavLink } from 'react-router-dom';
import './Cards.css'

const swipeVariants = {
  initial: { x: 0 },
  dragLeft: { x: -300, opacity: 0 },
  dragRight: { x: 300, opacity: 0 },
};

const Cards = (props) => {
  const { id } = useParams();
  const [cookies, setCookie] = useCookies(['AccessToken', 'UserId', 'LoadedCards']);
  const [cards, setCards] = useState([]);
  const [selectedCard, setSelectedCard] = useState(-1);
  const sessionId = id;

  const cardsFeedback = async (idx, direction) => {
    try {
      const response = await fetch(`http://localhost:8080/sessions/${sessionId}/scroll`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${cookies.AccessToken}`
          },
          body: JSON.stringify({
            'sessionID': sessionId,
            'placeID': idx,
            'is_liked': direction === "right" ? true : false
          })
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  const getCards = async () => {
    var response = await fetch('http://localhost:8080/sessions/'+ sessionId, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${cookies.AccessToken}`
        },
        body: JSON.stringify({
          'sessionID': sessionId
        })
    });
    var data = (await response.json()).session;
    const participant = data.users.find(participant => participant.ID === Number(cookies.UserId));
    response = await fetch('http://localhost:8080/cards', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${cookies.AccessToken}`
        },
        body: JSON.stringify({
          "categories": participant.Categories,
          "prompt" : participant.Request,
          'sessionID': sessionId
        })
    });
    data = await response.json();
    console.log(data.cards);
    setCards(data.cards ?? []);
  };

  useEffect(() => {

    const cardsFeedback = async () => {
      try {
        const response = await fetch(`http://localhost:8080/sessions/${sessionId}/check_match`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${cookies.AccessToken}`
            },
            body: JSON.stringify({
              'sessionID': sessionId
            })
        });
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
  
        const data = await response.json();
        if (data.is_matched) {
          window.location.reload();
        }
    
      } catch (error) {
        console.error('Error creating session:', error);
      }
    };

    getCards();
    const pollingInterval = setInterval(cardsFeedback, 1000);
    return () => clearInterval(pollingInterval);
   }, [cookies, sessionId]);

  const [xOffset, setXOffset] = useState(0);

  const handleDrag = (event, info) => {
    setXOffset(info.offset.x);
  };

  const handleDragEnd = () => {
    if (xOffset < -200) {
      swipeCard('left');
    } else if (xOffset > 200) {
      swipeCard('right');
    }
    setXOffset(0);
  };

  const swipeCard = (direction) => {
    cardsFeedback(cards[0].idx, direction);
    console.log(cards[0].idx, cards[0].title, direction);
    setCards(cards.slice(1));
    if (cards.length <= 1) {
      getCards();
    }
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

  useEffect(() => {
    var objects = document.getElementsByClassName('cards-body');
    Array.from(objects).map((item) => {
      const img = new Image();
      img.src = item.getAttribute('data-src');
      img.onload = () => {
        sessionStorage.setItem("LoadedCards", true);
        item.style.backgroundImage = `url(${item.getAttribute('data-src')})`;
      };
      img.onerror = () => {
        sessionStorage.setItem("LoadedCards", false);
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

  const handleSwipeOrClick = (index) => {
    if (xOffset === 0) {
      handleCardClick(index);
    }
  };

  const handleCardClick = (index) => {
    console.log(index);
    setSelectedCard(selectedCard === index ? -1 : index);
  };

  return (
    <div className={sessionStorage.getItem("LoadedCards")  ? "cards-body loadedCards" : "cards-body"} data-src="/bg_cards.png">
      <img src="data:image/gif;base64,R0lGODlhMgAbAIAAAP///wAAACH5BAEAAAEALAAAAAAyABsAAAIjjI+py+0Po5y02ouz3rz7D4biSJbmiabqyrbuC8fyTNf2zRUAOw==" data-src="/logo.png" class="logo" alt="Your Logo"></img>
      <ProfileHeader />
      <div class="cards-desc">
        <p>Листни карточку вправо, если место тебе понравилось, в противном случае - влево.</p>
      </div>
      <div style={{height: "80vh", display: "flex", alignItems: "center", justifyContent: "center"}}>
      {cards.slice().reverse().map((card, index) => (
        <motion.div
          key={index}
          onClick={() => handleSwipeOrClick(index)}
          animate={{ scale: selectedCard === index ? 1.15 : 1 }}
          drag="x"
          dragConstraints={{ top: 0, bottom: 0, left: 0, right: 0 }}
          dragElastic={0.8}
          dragMomentum={false}
          variants={swipeVariants}
          initial="initial"
          onDrag={handleDrag}
          onDragEnd={handleDragEnd}
          class="cards-card"
          whileTap={{ scale: 1.05 }}
          transition={{ duration: 0.3, type: 'spring', stiffness: 300 }}
          style={{ boxShadow: xOffset < -30 ? "0 0 20px red" : xOffset > 30 ? "0 0 20px green" : "none" }}
        >
          <div style={{display: "flex", flexDirection: "column", justifyContent: "space-between", }}>
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
      ))}
      </div>
    </div>
  );
};

export default Cards;