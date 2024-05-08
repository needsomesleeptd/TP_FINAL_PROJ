import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import { NavLink } from 'react-router-dom';
import './Cards.css'
import './Main.css'

const Match = (props) => {
 const { id } = useParams();
 const [cookies] = useCookies(['AccessToken', 'UserId']);
 const [cards, setCards] = useState([]);
 const sessionId = id;

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
      console.log(data.cards);
      setCards(data.cards);
  
    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  cardsFeedback();
 }, [cookies, sessionId]);


 return (
  <div class="cards-body">
    <div class="spec">
     <div style={{ textAlign: "center", marginTop: 50 }}>
       <p style={{ fontSize: "30px", fontWeight: "bold", color: "white" }}>Мы нашли подходящее место для вас!</p>
     </div>
     <div style={{ height: "80vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
       {cards.slice().reverse().map((card, index) => (
         <div
           key={index}
           class="cards-card"
           style={{marginTop: -20}}
         >
           <div style={{display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center"}}>
             <img src={card.image} alt="" class="cards-img" />
             <p style={{textAlign: "center" }}>{card.title}</p>
           </div>
         </div>
       ))}
     </div>
     <div style={{ textAlign: "center", marginTop: -100 }}>
     <NavLink to="/" style={{ textDecoration: "none" }}>
        <button class="modal-button" style={{marginBottom: "20px"}}>Вернуться</button>
      </NavLink>
    </div>
   </div>
   </div>
 );
};

export default Match;
