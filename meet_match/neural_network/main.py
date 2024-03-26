from flask import Flask, jsonify, request
import numpy as np
from model import train, predict, load

vectorizer, vectorized_matrix = None, None

app = Flask(__name__)


@app.route("/rec", methods=["POST", "GET"])
def image_post_request():
    query = request.json
    # print(query)
    res_data = predict(
        "../database/dist.csv",
        vectorizer,
        vectorized_matrix,
        query["query"],
        query["from_line"],
        query["to_line"],
    )
    # print(res_data)
    return jsonify(res_data)


if __name__ == "__main__":
    vectorizer, vectorized_matrix = load("vectorizer.pkl", "vectorized_matrix.pkl")
    app.run(host="0.0.0.0", port=5000)
