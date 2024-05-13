from flask import Flask, jsonify, request
import numpy as np
from helpers import load
import pickle
from model import RecommendationSystem

import pandas as pd
import json

model = RecommendationSystem()
app = Flask(__name__)
model.build_embeddings()


@app.route("/rec", methods=["POST", "GET"])
def image_post_request():
    query = request.json
    recs = model.get_rec(query["user_id"], query["session_id"], query["query"], num_idx=10)
    return jsonify({"recs": recs})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
