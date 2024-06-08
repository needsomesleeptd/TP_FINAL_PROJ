import requests
from bs4 import BeautifulSoup
import csv


url = "https://kudago.com/public-api/v1.4/events/?location=msk&categories=cinema,concert,education,entertainment,exhibition,fashion,festival,holiday,party,quest,recreation,theater,tour&fields=id,dates,title,short_title,slug,place,description,body_text,location,categories,tagline,age_restriction,price,is_free,images,favorites_count,comments_count,site_url,tags,participants&actual_since=1715621332&actual_until=1716215233"


def collect_data(url):
    resp = requests.get(url)
    resp_json = resp.json()

    with open("events.csv", "w", encoding="utf-8") as dst:
        fieldnames = [
            "id",
            "dates",
            "title",
            "short_title",
            "slug",
            "place",
            "description",
            "body_text",
            "location",
            "categories",
            "tagline",
            "age_restriction",
            "price",
            "is_free",
            "images",
            "favorites_count",
            "comments_count",
            "site_url",
            "tags",
            "participants",
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
                        "dates": i["dates"],
                        "title": i["title"],
                        "short_title": i["short_title"],
                        "slug": i["slug"],
                        "place": i["place"],
                        "description": i["description"],
                        "body_text": i["body_text"],
                        "location": i["location"],
                        "categories": i["categories"],
                        "tagline": i["tagline"],
                        "age_restriction": i["age_restriction"],
                        "price": i["price"],
                        "is_free": i["is_free"],
                        "images": i["images"],
                        "favorites_count": i["favorites_count"],
                        "comments_count": i["comments_count"],
                        "site_url": i["site_url"],
                        "tags": i["tags"],
                        "participants": i["participants"],
                    }
                )
            resp = requests.get(resp_json["next"])
            resp_json = resp.json()


if __name__ == "__main__":
    collect_data(url)
