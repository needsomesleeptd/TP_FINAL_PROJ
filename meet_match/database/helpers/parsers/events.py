import requests
from bs4 import BeautifulSoup
import csv


url = "https://kudago.com/public-api/v1.4/places/?location=msk&categories=amusement,anticafe,animal-shelters,art-centers,art-space,attractions,bar,brewery,cats,cinema,clubs,concert-hall,comedy-club,dance-studio,dogs,homesteads,handmade,museums,park,questroom,restaurants,sights,theatre&fields=id,title,short_title,address,location,timetable,phone,images,description,body_text,foreign_url,subway,coords,favorites_count,comments_count,is_closed,tags,categories"


def collect_data(url):
    resp = requests.get(url)
    resp_json = resp.json()

    with open("dist.csv", "w", encoding="utf-8") as dst:
        fieldnames = [
            "id",
            "title",
            "short_title",
            "address",
            "location",
            "timetable",
            "phone",
            "images",
            "description",
            "body_text",
            "foreign_url",
            "subway",
            "coords",
            "favorites_count",
            "comments_count",
            "is_closed",
            "tags",
            "categories",
        ]

        writer = csv.DictWriter(dst, fieldnames=fieldnames)
        writer.writeheader()

        counter = 0
        while resp_json["next"] and counter < 50:
            counter += 1
            for i in resp_json["results"]:
                writer.writerow(
                    {
                        "id": i["id"],
                        "title": i["title"],
                        "short_title": i["short_title"],
                        "address": i["address"],
                        "location": i["location"],
                        "timetable": i["timetable"],
                        "phone": i["phone"],
                        "images": i["images"],
                        "description": i["description"],
                        "body_text": i["body_text"],
                        "foreign_url": i["foreign_url"],
                        "subway": i["subway"],
                        "coords": i["coords"],
                        "favorites_count": i["favorites_count"],
                        "comments_count": i["comments_count"],
                        "is_closed": i["is_closed"],
                        "tags": i["tags"],
                        "categories": i["categories"],
                    }
                )
            resp = requests.get(resp_json["next"])
            resp_json = resp.json()


if __name__ == "__main__":
    collect_data(url)
