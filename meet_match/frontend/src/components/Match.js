import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useCookies } from 'react-cookie';

const Match = (props) => {
 const { id } = useParams();
 const [cookies] = useCookies(['AccessToken', 'UserId']);
 const [cards, setCards] = useState([]);
 const sessionId = id;

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
      console.log(data.cards);
      setCards(data.cards);
  
    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  cardsFeedback();
 }, [cookies, sessionId]);


 return (
   <div>
     <div style={{ textAlign: "center", marginBottom: 20 }}>
       <p style={{ fontSize: "24px", fontWeight: "bold" }}>Мы нашли подходящее место для вас!</p>
     </div>
     <div style={{ height: "80vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
       {cards.slice().reverse().map((card, index) => (
         <div
           key={index}
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
         </div>
       ))}
     </div>
     <div style={{ textAlign: "center", marginTop: 20 }}>
     <a href="/" style={{ textDecoration: "none" }}>
        <button style={{ marginBottom: "40px", padding: "10px 20px", fontSize: "16px", backgroundColor: "lightblue", border: "2px solid #000", borderRadius: "5px", cursor: "pointer" }}>Вернуться на главную страницу</button>
      </a>
    </div>
   </div>
 );
};

export default Match;
