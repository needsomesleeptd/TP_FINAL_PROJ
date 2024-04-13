import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useCookies } from 'react-cookie';

const swipeVariants = {
  initial: { x: 0 },
  dragLeft: { x: -300, opacity: 0 },
  dragRight: { x: 300, opacity: 0 },
};

const Cards = () => {
  const [cookies, setCookie] = useCookies(['meetmatchname', 'meetmatchsession', 'meetmatchrequest']);
  const [cards, setCards] = useState([]);

  const cardsFeedback = async (direction) => {
    try {
      const response = await fetch('http://localhost:8080/cards', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('jwt')}`
          },
          body: JSON.stringify({
            'sessionID': cookies.meetmatchsession,
            'userIDToModify': Number(localStorage.getItem('userID')),
            'placeID': 0,
            'isLiked': direction === "right" ? true : false
          })
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  useEffect(() => {
    const getCards = async () => {
      var response = await fetch('http://localhost:8080/sessions/'+ cookies.meetmatchsession, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('jwt')}`
          },
          body: JSON.stringify({
            'sessionID': cookies.meetmatchsession
          })
      });
      var data = await response.json();
      response = await fetch('http://localhost:8080/cards', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('jwt')}`
          },
          body: JSON.stringify({
            "prompt" : cookies.meetmatchrequest,
            "page" : 1,
            "cardsPerPage" : 10
          })
      });
      data = await response.json();
      console.log(data.cards);
      setCards(data.cards);
    };

    getCards();
  }, []);

  const [swipedCard, setSwipedCard] = useState(null);
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
    cardsFeedback(direction);
    console.log(cards);
    console.log(swipeCard);
    setSwipedCard(cards[0]);
    setCards(cards.slice(1));
  };

  return (
    <div>
      <div style={{ textAlign: "center", marginBottom: 20 }}>
        <p style={{ fontSize: "24px", fontWeight: "bold" }}>Свайпай карточки влево или вправо</p>
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
          animate={index + 1 != cards.length ? "initial" : xOffset < -200 ? "dragLeft" : xOffset > 200 ? "dragRight" : "initial"}
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
            <img src={card.image} style={{
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