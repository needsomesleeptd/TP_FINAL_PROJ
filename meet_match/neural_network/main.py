from flask import Flask, jsonify, request
import numpy as np
from model import train,predict


app = Flask(__name__)

@app.route('/rec', methods=['POST'])
def image_post_request():  
    query = request.json['query']
    res_data = predict('vectorizer.pkl', 'vectorized_matrix.pkl', query)
    return jsonify(res_data)




if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5000)