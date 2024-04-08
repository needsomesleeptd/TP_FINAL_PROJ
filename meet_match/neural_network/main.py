from flask import Flask, jsonify, request
import numpy as np
from helpers import load
import pickle
from model import predict

import pandas as pd
import json

text_embeddings = None

app = Flask(__name__)


class ModelRequest:
    def __init__(
        self,
        query: str = None,
        label: str = None,
        from_line: int = 0,
        to_line: int = -1,
    ) -> None:
        self.query = query
        self.label = label
        self.from_line = from_line
        self.to_line = to_line


# type Card struct {
# 	ImgUrl   string `json:"image"`
# 	CardName string `json:"title,card_name"`
# 	Rating   int    `json:"rating,omitempty"`
# }


class CardResponse:
    def __init__(self, imgurl, cardname, rating=-1) -> None:
        self.image = imgurl
        self.card_name = cardname
        self.rating = rating

    def __str__(self) -> str:
        return f"image = {self.image}, card_name = {self.card_name}, rating = {self.rating}"


def get_cards(csv_name, res_data):
    cards = []
    df = pd.read_csv(csv_name)
    for i in res_data[0]:
        cards.append(
            {
                "idx": i[0],
                "title": df["title"][i[0]],
                "image": json.loads(df["images"][i[0]].replace("'", '"'))[0]["image"],
                "cosine_similarities": i[1],
            }
        )

    return cards


@app.route("/rec", methods=["POST", "GET"])
def image_post_request():
    query = request.json
    # print(query)
    mr = ModelRequest(**query)
    res_data = predict(
        "../database/dist.csv",
        text_embeddings,
        mr.query,
        mr.from_line,
        mr.to_line,
    )
    res = get_cards("../database/dist.csv", res_data)
    return jsonify(res)


if __name__ == "__main__":
    with open("data_e5_vectors.pkl", "rb") as f:
        text_embeddings = pickle.load(f)
    app.run(host="0.0.0.0", port=5000)
