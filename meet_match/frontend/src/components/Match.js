import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import { motion } from 'framer-motion';
import { NavLink } from 'react-router-dom';
import './Cards.css'
import './Main.css'

const Match = (props) => {
 const { id } = useParams();
 const [cookies] = useCookies(['AccessToken', 'UserId']);
 const [cards, setCards] = useState([]);
 const [selectedCard, setSelectedCard] = useState(-1);
 const sessionId = id;

 useEffect(() => {
  const cardsFeedback = async () => {
    try {
      const response = await fetch(`http://localhost:8080/sessions/${sessionId}/matches`, {
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
      console.log(data.cards);
      setCards(data.cards);
  
    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  cardsFeedback();
 }, [cookies, sessionId]);

 const cardsContinue = async () => {
  try {
    const response = await fetch(`http://localhost:8080/sessions/${sessionId}/continueScrolling`, {
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

    window.location.reload();

  } catch (error) {
    console.error('Error creating session:', error);
  }
};


const handleCardClick = (index) => {
  console.log(index);
  setSelectedCard(selectedCard === index ? -1 : index);
};


 return (
  <div class="cards-body">
    <div class="spec">
     <div style={{ textAlign: "center", marginTop: 20 }}>
      {
        cards.length > 1 ?
        <p style={{ fontSize: "30px", fontWeight: "bold", color: "white" }}>Мы нашли подходящие места для вас!</p>
        :
        <p style={{ fontSize: "30px", fontWeight: "bold", color: "white" }}>Мы нашли подходящее место для вас!</p>
      }
     </div>
     {
      cards.length > 2 ?
      <div className="horizontal-scroll-block">
        {cards.slice().reverse().map((card, index) => (
          <div
            key={index}
            class="cards-card"
            onClick={() => handleCardClick(index)}
            style={{ position: "relative", marginBottom: "20px" }}
          >
           <div style={{display: "flex", flexDirection: "column", justifyContent: "space-between", width: "270px" }}>
            {selectedCard === index ? (
              <>
                <p style={{textAlign: "center", margin: "20px 10px", fontSize: "16px", fontWeight: "bold" }}>{card.title}</p>
                <p className="cards-descr" dangerouslySetInnerHTML={{ __html: card.description }} />
                {card.age_restriction && <p style={{marginLeft: "10px", textAlign: "left", marginBottom: "10px", fontSize: "13px" }}>Возраст: {card.age_restriction}</p>}
                {card.cost && <p style={{marginLeft: "10px", textAlign: "left", marginBottom: "10px", fontSize: "13px" }}>Цена: {card.cost}</p>}
                {card.timetable && <p style={{marginLeft: "10px", textAlign: "left", marginBottom: "10px", fontSize: "13px" }}>Расписание: {card.timetable}</p>}
                {card.subway && <p style={{marginLeft: "10px", textAlign: "left", marginBottom: "10px", fontSize: "13px" }}>Метро: {card.subway}</p>}
                {card.site_url && <p style={{marginLeft: "10px",  textAlign: "left", marginBottom: "10px", fontSize: "13px" }}>Сайт: <a href={card.site_url} target="_blank" rel="noopener noreferrer"
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
          </div>
        ))}
      </div>
      :
      <div className="horizontal-scroll-block" style={{justifyContent: "center"}}>
        {cards.slice().reverse().map((card, index) => (
          <div
            key={index}
            class="cards-card"
            onClick={() => handleCardClick(index)}
            style={{ position: "relative", marginBottom: "20px" }}
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
          </div>
        ))}
      </div>
     }
     
     <div style={{ textAlign: "center", marginTop: "30px" }}>
      <div>
      <button onClick={cardsContinue} class="modal-button" style={{marginBottom: "10px"}}>Продолжить поиск</button>
      </div>
     <NavLink to="/" style={{ textDecoration: "none" }}>
        <button class="modal-button" style={{marginBottom: "10px"}}>На главную</button>
      </NavLink>
    </div>
   </div>
   </div>
 );
};

export default Match;
