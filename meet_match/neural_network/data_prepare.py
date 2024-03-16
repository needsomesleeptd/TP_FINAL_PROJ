import re
import pymorphy3


lemmatizer = pymorphy3.MorphAnalyzer()
lemmatizer_cache = {}

def lemmatize_token(token):
    if lemmatizer.word_is_known(token):
        if token not in lemmatizer_cache:
            lemmatizer_cache[token] = lemmatizer.parse(token)[0].normal_form
        return lemmatizer_cache[token]
    return token

def preprocess_text(text, token_pattern):
    tokens = re.findall(token_pattern, str(text).lower())
    lemmatized_tokens = [lemmatize_token(token) for token in tokens]
    return " ".join(lemmatized_tokens)

def preprocess_all_texts(texts, token_pattern):
    preprocessed_texts = []
    for text in texts:
        preprocessed_text = preprocess_text(text, token_pattern)
        preprocessed_texts.append(preprocessed_text)
    return preprocessed_texts
