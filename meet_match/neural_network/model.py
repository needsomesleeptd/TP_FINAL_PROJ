import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity 
import pandas as pd
import pickle
from data_prepare import preprocess_all_texts, preprocess_text

TOKEN_PATTERN = "[a-zA-Zа-яА-ЯёЁ]+"

def train(csv_file, vectorizer_file, vectorized_matrix_file):
    df = pd.read_csv(csv_file)
    df_descr = df['tags'].fillna('') + '         ' + df['description'].fillna('') + '         ' + df['body_text'].fillna('')
    descr_lemm = preprocess_all_texts(df_descr, TOKEN_PATTERN)
    vectorizer = TfidfVectorizer(min_df=2, norm=None, ngram_range=(1,5))
    vectorized_matrix = vectorizer.fit_transform(descr_lemm)
    with open(vectorizer_file, 'wb') as file:
        pickle.dump(vectorizer, file)
    with open(vectorized_matrix_file, 'wb') as file:
        pickle.dump(vectorized_matrix, file)

def predict(vectorizer_file, vectorized_matrix_file, sample_quary):
    with open(vectorizer_file, 'rb') as file:
        vectorizer = pickle.load(file)
    with open( vectorized_matrix_file, 'rb') as file:
        vectorized_matrix = pickle.load(file)
    sample_query = preprocess_text(sample_quary, TOKEN_PATTERN)
    qery_tdidf = vectorizer.transform([sample_query])
    cosine_similarities = cosine_similarity(qery_tdidf, vectorized_matrix)
    top_indices = cosine_similarities.argsort()[0][-5:][::-1]
    res = []
    for idx in top_indices:
        res.append({
            "idx": idx,
            "cosine_similarities": cosine_similarities[0][idx]
        })
    return res

# train('dist.csv', 'vectorizer.pkl', 'vectorized_matrix.pkl')
