package models

type Card struct {
	Id             uint64
	ImgUrl         string
	CardName       string
	Rating         *uint64 //  star fields can be nil
	Description    *string
	Subway         *string
	Cost           *string
	Timetable      *string
	AgeRestriction *string
	Phone          *string
	SiteUrl        *string
}
