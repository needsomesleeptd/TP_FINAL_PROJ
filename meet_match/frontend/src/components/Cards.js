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
  const [cookies] = useCookies(['AccessToken', 'UserId']);
  const [cards, setCards] = useState([]);
  const sessionId = id;

  const cardsFeedback = async (idx, direction) => {
    try {
      const response = await fetch(`/api/sessions/${sessionId}/scroll`, {
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
    var response = await fetch('/api/sessions/'+ sessionId, {
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
    response = await fetch('/api/cards', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${cookies.AccessToken}`
        },
        body: JSON.stringify({
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
        const response = await fetch(`/api/sessions/${sessionId}/check_match`, {
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
    const pollingInterval = setInterval(cardsFeedback, 500);
    return () => clearInterval(pollingInterval);
   }, [cookies, sessionId]);

  const [xOffset, setXOffset] = useState(0);

  const handleDrag = (event, info) => {
    setXOffset(info.offset.x);
  };

  const handleDragEnd = () => {
    setXOffset(0);

    if (xOffset < -200) {
      swipeCard('left');
    } else if (xOffset > 200) {
      swipeCard('right');
    }
  };

  const swipeCard = (direction) => {
    cardsFeedback(cards[0].id, direction);
    console.log(cards[0].id, direction);
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

  return (
    <div class="cards-body">
      <img src="/logo.png" class="logo" alt="Your Logo"></img>
      <ProfileHeader />
      <div class="cards-desc">
        <p>Давай подберём подходящее место для вашей встречи. Для этого просто листни карточку</p>
        <p>влево, если место тебе не нравится, или же вправо, если место тебе понравилось.</p>
      </div>
      <div style={{height: "80vh", display: "flex", alignItems: "center", justifyContent: "center"}}>
      {cards.slice().reverse().map((card, index) => (
        <motion.div
          key={index}
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
          <div style={{display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center"}}>
            <img src={card.image} alt="" class="cards-img" />
            <p style={{textAlign: "center" }}>{card.title}</p>
          </div>
        </motion.div>
      ))}
      </div>
    </div>
  );
};

export default Cards;