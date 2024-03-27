import React, { useState } from 'react';
import { motion } from 'framer-motion';

const swipeVariants = {
  initial: { x: 0 },
  dragLeft: { x: -300, opacity: 0 },
  dragRight: { x: 300, opacity: 0 },
};

const Cards = () => {
  const [cards, setCards] = useState([
    { id: 1, imageUrl: 'https://avatars.mds.yandex.net/i?id=8baa27866533b4b9ad9cd4e3d7bde320b7d9d298-10411335-images-thumbs&n=13', caption: 'Школьный снюсоед' },
    { id: 2, imageUrl: 'https://avatars.mds.yandex.net/i?id=076385ea095bdedb124098355087d0c6429db46f-4409557-images-thumbs&n=13', caption: 'Забивной снюсоед' },
    { id: 3, imageUrl: 'https://www.meme-arsenal.com/memes/b6f8b1ec5533277508251f186163454e.jpg', caption: 'Домашний снюсоед' },
  ]);

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
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            position: "absolute",
          }}
        >
          <div style={{display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center"}}>
            <img src={card.imageUrl} style={{
                marginTop: "30px",
                width: "200px",
                height: "200px",
                objectFit: "cover",
                borderRadius: "10%",
                pointerEvents: "none"
                }} />
            <p style={{textAlign: "center" }}>{card.caption}</p>
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