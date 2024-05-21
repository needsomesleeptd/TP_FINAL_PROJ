package models

type CardStats struct {
	CardID      uint64
	SwipedTimes uint64
}

type PersonScrollStats struct {
	PersonalStats        PersonalScrollStats
	SessionsCount        uint64
	MostDislikedPlace    Card
	MostLikedPlace       Card
	MostlikedScrolled    uint64
	MostDislikedScrolled uint64
}

type PersonalScrollStats struct {
	Swiped         uint64
	PoisitveSwipes uint64
	NegativeSwipes uint64
}
