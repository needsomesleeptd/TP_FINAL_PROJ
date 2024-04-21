import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useParams } from 'react-router-dom';
import { useCookies } from 'react-cookie';

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
    cardsFeedback(cards[0].idx, direction);
    console.log(cards[0].idx, direction);
    setCards(cards.slice(1));
    if (cards.length == 0) {
      getCards();
    }
  };

  return (
    <div>
      <div style={{ textAlign: "center", marginBottom: 20 }}>
        <p style={{ fontSize: "24px", fontWeight: "bold" }}>Давай подберём подходящее место для вашей встречи</p>
        <p style={{ fontSize: "24px", fontWeight: "bold" }}>Для этого просто листни карточку влево, если место тебе не нравится</p>
        <p style={{ fontSize: "24px", fontWeight: "bold" }}>Или же листни карточку вправо, если место тебе понравилось</p>
      </div>
      <div style={{height: "80vh", display: "flex", alignItems: "center", justifyContent: "center"}}>
      {cards.slice().reverse().map((card, index) => (
        <motion.div
          key={index}
          drag
          dragConstraints={{ top: 0, bottom: 0, left: 0, right: 0 }}
          dragElastic={0.8}
          dragMomentum={false}
          variants={swipeVariants}
          initial="initial"
          animate={index + 1 !== cards.length ? "initial" : xOffset < -200 ? "dragLeft" : xOffset > 200 ? "dragRight" : "initial"}
          onDrag={handleDrag}
          onDragEnd={handleDragEnd}
          style={{
            width: 300,
            height: 400,
            background: "lightblue",
            borderRadius: 20,
            border: "2px solid #000",
            boxSizing: "border-box",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            position: "absolute",
          }}
        >
          <div style={{display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center"}}>
            <img src={card.image} alt="" style={{
                marginTop: "30px",
                width: "200px",
                height: "200px",
                border: "2px solid #000",
                boxSizing: "border-box",
                objectFit: "cover",
                borderRadius: "10%",
                pointerEvents: "none"
                }} />
            <p style={{textAlign: "center" }}>{card.title}</p>
          </div>
          <p style={{ position: "absolute", top: 0, left: 20 }}>Не нравится</p>
          <p style={{ position: "absolute", top: 0, right: 20 }}>Нравится</p>
        </motion.div>
      ))}
      </div>
    </div>
  );
};

export default Cards;