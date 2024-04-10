import pickle


def load(vectorizer_file, vectorized_matrix_file):
    with open(vectorizer_file, "rb") as file:
        vectorizer = pickle.load(file)
    with open(vectorized_matrix_file, "rb") as file:
        vectorized_matrix = pickle.load(file)
    return vectorizer, vectorized_matrix
