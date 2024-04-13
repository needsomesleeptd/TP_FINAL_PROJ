# нужно переписать _get_swipes_for_session и _load_embeddings

'''
пример работы:
rec_sys =RecommendationSystem()
rec_sys.get_rec(1,1,"яблоко", 10)
'''


import torch
import pandas as pd
import torch.nn.functional as F
import pickle
import sqlite3
import random

from torch import Tensor
from transformers import AutoTokenizer, AutoModel
from tqdm import tqdm_notebook
from typing import List, Dict, Any, Tuple


class RecommendationSystem:
    def __init__(self):
        self.tokenizer = AutoTokenizer.from_pretrained(
            'intfloat/multilingual-e5-large-instruct')
        self.model = AutoModel.from_pretrained(
            'intfloat/multilingual-e5-large-instruct')
        self.db_path = 'prog_db.db'
        self.embed_dict = self._load_embeddings()

    def _get_swipes_for_session(self, user_id: int, session_id: int, for_group: bool = False) -> List[Tuple[int, bool]]:
        """
        Возвращает список пар (place_id, is_liked) для заданного пользователя и сессии.

        Параметры:
        - db_path (str): Путь к файлу базы данных SQLite.
        - user_id (int): Идентификатор пользователя.
        - session_id (int): Идентификатор сессии.
        - for_group (bool): Свайпы всей группы, без user_id, или наоборот

        Возвращает:
        - List[Tuple[int, bool]]: Список пар, где каждая пара содержит ID места и флаг, указывающий,
                                  понравилось ли место пользователю (True - понравилось, False - не понравилось).
        """

        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        query = ''

        if for_group:
            query = """
            SELECT place_id, is_liked
            FROM place_swipes
            WHERE user_id != ? AND session_id = ?
            """
        else:
            query = """
            SELECT place_id, is_liked
            FROM place_swipes
            WHERE user_id = ? AND session_id = ?
            """

        cursor.execute(query, (user_id, session_id))

        swipes = cursor.fetchall()
        conn.close()

        swipes_list = [(place_id, bool(is_liked))
                       for place_id, is_liked in swipes]

        return swipes_list

    def _load_embeddings(self) -> Dict[int, Any]:
        """
        Загружает эмбеддинги из базы данных SQLite и десериализует их.

        Параметры:
        - db_path (str): Путь к файлу базы данных SQLite.

        Возвращает:
        - Dict[int, Any]: Словарь, где ключом является ID места (place_id), а значением - десериализованный эмбеддинг.
        """

        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        cursor.execute("SELECT place_id, embedding FROM embeddings")
        embeddings_data = cursor.fetchall()
        conn.close()

        embeddings = {place_id: pickle.loads(
            embedding) for place_id, embedding in embeddings_data}
        return embeddings

    def _cosine_similarity(self, embeddings: Tensor, all_embeddings: Tensor) -> Tensor:
        """
        Рассчитывает косинусное сходство между двумя наборами векторов эмбеддингов.

        Параметры:
        - embeddings (Tensor): Тензор эмбеддингов размерностью (n, d), где n - количество эмбеддингов,
                               а d - размерность каждого эмбеддинга.
        - all_embeddings (Tensor): Тензор всех эмбеддингов размерностью (m, d), где m - количество эмбеддингов
                                   во втором наборе, а d - размерность каждого эмбеддинга.

        Возвращает:
        - Tensor: Тензор косинусного сходства размерностью (n, m), где каждый элемент [i, j] 
                  представляет собой косинусное сходство между i-тым эмбеддингом из первого набора 
                  и j-тым эмбеддингом из второго набора.
        """

        # Normalize the embeddings to unit vectors.
        # embeddings_norm = F.normalize(embeddings)
        # all_embeddings_norm = F.normalize(all_embeddings)

        return torch.mm(embeddings, all_embeddings.T)
        # return (embeddings @ all_embeddings.T)

    def __get_detailed_instruct(self, task_description: str, query: str) -> str:
        return f'Instruct: {task_description}\nQuery: {query}'

    def __average_pool(self, last_hidden_states: Tensor,
                       attention_mask: Tensor) -> Tensor:
        last_hidden = last_hidden_states.masked_fill(
            ~attention_mask[..., None].bool(), 0.0)
        return last_hidden.sum(dim=1) / attention_mask.sum(dim=1)[..., None]

    def _generate_embedding(self, text):
        inputs = self.tokenizer(text, return_tensors='pt',
                                padding=True, truncation=True, max_length=512)
        with torch.no_grad():
            outputs = self.model(**inputs)

        # embedding = outputs.last_hidden_state.mean(dim=1).numpy()
        pooled_embeddings = self.__average_pool(
            outputs.last_hidden_state, inputs['attention_mask'])
        embedding_normalized = F.normalize(pooled_embeddings, p=2, dim=1)
        return embedding_normalized

    def _get_rec_on_query(self, query: str, used_ids: List[int], num_items: int = 20) -> List[int]:
        """
        Получает рекомендации на основе текстового запроса, исключая уже использованные ID.

        Параметры:
        - query (str): Текстовый запрос для поиска.
        - embeddings (Dict[int, torch.Tensor]): Словарь с эмбеддингами мест, где ключ - это ID места.
        - used_ids (List[int]): Список ID мест, которые уже были использованы и должны быть исключены из результатов.
        - num_items (int): Количество возвращаемых рекомендаций.

        Возвращает:
        - List[int]: Список ID мест, которые рекомендуются на основе запроса.
        """

        task = 'Represent this sentence for searching relevant passages: '
        query_embedding = self._generate_embedding(
            self.__get_detailed_instruct(task, query))

        place_ids = list(self.embed_dict.keys())
        embedding_matrix = torch.cat(list(self.embed_dict.values()), dim=0)

        similarities = self._cosine_similarity(
            query_embedding, embedding_matrix).squeeze(0)
        sorted_scores, sorted_indices = torch.sort(
            similarities, descending=True)

        top_n_place_ids = []
        idx_iter = 0
        while len(top_n_place_ids) < num_items and idx_iter < len(place_ids):
            current_idx = sorted_indices[idx_iter].item()

            if place_ids[current_idx] not in used_ids:
                top_n_place_ids.append(place_ids[current_idx])

            idx_iter += 1

        return top_n_place_ids

    def _generate_rec_on_user_session_hist(
        self,
        user_id: int,
        session_id: int,
        swiped_places: List[Tuple[int, bool]],
        used_ids: List[int],
        num_top_idx: int = 20,
        dislike_weight: float = 0.5
    ) -> List[int]:
        """
        Генерирует рекомендации на основе истории лайков и дизлайков пользователя в данной сессии.

        Параметры:
        - user_id (int): Идентификатор пользователя.
        - session_id (int): Идентификатор сессии.
        - text_embeddings (Dict[int, torch.Tensor]): Словарь с эмбеддингами мест.
        - swiped_places: List[Tuple[int, bool]] : Список просвайпанных id, где свайп влево это 1.
        - used_ids (List[int]): Список ID мест, которые уже были использованы.
        - num_top_idx (int): Количество мест для возвращения.
        - dislike_weight (float): Вес, с которым учитываются дизлайки при вычислении рекомендаций.

        Возвращает:
        - List[int]: Список ID рекомендованных мест.
        """
        if not swiped_places:
            return []

        # Получение ID мест, которые пользователь лайкнул и дизлайкнул в текущей сессии
        liked_place_ids = [idx for (idx, is_liked)
                           in swiped_places if is_liked == 1]
        disliked_place_ids = [
            idx for (idx, is_liked) in swiped_places if is_liked == 0]

        # Извлечение эмбеддингов для лайкнутых и дизлайкнутых мест
        liked_embeddings = torch.cat(
            [self.embed_dict[i] for i in liked_place_ids], dim=0)
        disliked_embeddings = torch.cat(
            [self.embed_dict[i] for i in disliked_place_ids], dim=0)

        # Подготовка общей матрицы эмбеддингов и списка ID мест
        all_embeddings = torch.cat(list(self.embed_dict.values()), dim=0)
        place_ids = list(self.embed_dict.keys())

        # Расчёт косинусного сходства для лайкнутых и дизлайкнутых мест
        liked_similarities = self._cosine_similarity(
            liked_embeddings, all_embeddings)
        disliked_similarities = self._cosine_similarity(
            disliked_embeddings, all_embeddings)

        # Комбинирование сходства, учитывая вес дизлайков
        combined_scores = liked_similarities.mean(
            dim=0) - dislike_weight * disliked_similarities.mean(dim=0)

        # Сортировка мест на основе комбинированного сходства
        scores, indices = combined_scores.squeeze(0).sort(descending=True)

        # Формирование списка рекомендованных ID мест, исключая уже использованные
        top_n_place_ids = []
        idx_iter = 0
        while len(top_n_place_ids) < num_top_idx and idx_iter < len(indices):
            current_idx = indices[idx_iter].item()
            if place_ids[current_idx] not in used_ids:
                top_n_place_ids.append(place_ids[current_idx])
            idx_iter += 1

        return top_n_place_ids

    def _generate_group_rec_on_likes(
        self,
        user_id: int,
        session_id: int,
        swiped_places: List[Tuple[int, bool]],
        used_ids: List[int],
        num_top_idx: int = 20
    ) -> List[int]:
        """
        Генерирует рекомендации на основе лайков группы пользователей в данной сессии.

        Параметры:
        - user_id (int): Идентификатор пользователя.
        - session_id (int): Идентификатор сессии.
         - swiped_places: List[Tuple[int, bool]] : Список просвайпанных id, где свайп влево это 1.
        - used_ids (List[int]): Список ID мест, которые уже были использованы.
        - num_top_idx (int): Количество мест для возвращения.

        Возвращает:
        - List[int]: Список ID рекомендованных мест.
        """
        if not swiped_places:
            return []
        group_liked_place_ids = [
            idx for (idx, is_liked) in swiped_places if is_liked == 1]

        if not group_liked_place_ids:
            return []
        group_liked_embeddings = torch.cat(
            [self.embed_dict[i] for i in group_liked_place_ids], dim=0)

        all_embeddings = torch.cat(list(self.embed_dict.values()), dim=0)
        place_ids = list(self.embed_dict.keys())

        group_liked_similarities = self._cosine_similarity(
            group_liked_embeddings, all_embeddings)

        combined_scores = group_liked_similarities.mean(dim=0)

        scores, indices = combined_scores.squeeze(0).sort(descending=True)

        top_n_place_ids = []
        idx_iter = 0
        while len(top_n_place_ids) < num_top_idx and idx_iter < len(indices):
            current_idx = indices[idx_iter].item()
            if place_ids[current_idx] not in used_ids:
                top_n_place_ids.append(place_ids[current_idx])
            idx_iter += 1

        return top_n_place_ids

    def _get_random_places(self, all_place_ids, used_indices, num_places=10):
        place_ids = []
        idx_iter = 0

        while len(place_ids) < num_places and idx_iter < len(all_place_ids):
            idx = random.randint(0, len(all_place_ids))
            if all_place_ids[idx] not in used_indices:
                place_ids.append(all_place_ids[idx])
            idx_iter += 1
        return place_ids

    def get_rec(self, user_id: int, session_id: int, query: str, num_idx=50) -> List[int]:
        """   
        Параметры:
        - user_id (int): Идентификатор пользователя.
        - session_id (int): Идентификатор сессии.
        - query (str): Текстовый запрос для поиска.
        - num_idx (int) : Число индексов на выходе.

        Возвращает:
        - List[int]: Список ID рекомендованных мест.
        """

        proportions = [0.3, 0.3, 0.3, 0.1]  # лучше динамически подбирать, пока так(
        parts = [int(num_idx * p) for p in proportions]

        recomendation_indices = []
        swipe_user_hist = self._get_swipes_for_session(
            user_id, session_id, for_group=False)
        used_indices = swipe_user_hist.copy()

        query_rec_list = self._get_rec_on_query(query, used_indices, parts[0])
        if len(query_rec_list) < parts[0]:
            query_rec_list = query_rec_list + self._get_random_places(
                list(self.embed_dict.keys()), used_indices, len(query_rec_list) - parts[0])

        used_indices = used_indices + query_rec_list

        rec_user_hist = self._generate_rec_on_user_session_hist(
            user_id, session_id, swipe_user_hist, used_indices, parts[1])
        if len(rec_user_hist) < parts[1]:
            rec_user_hist = rec_user_hist + self._get_random_places(
                list(self.embed_dict.keys()), used_indices, len(rec_user_hist) - parts[1])
        used_indices = used_indices + rec_user_hist

        swipe_group_hist = self._get_swipes_for_session(
            user_id, session_id, for_group=True)

        rec_group_hist = self._generate_group_rec_on_likes(
            user_id, session_id, swipe_group_hist, used_indices, parts[2])
        if len(rec_group_hist) < parts[2]:
            rec_group_hist = rec_group_hist + self._get_random_places(
                list(self.embed_dict.keys()), used_indices, len(rec_group_hist) - parts[2])
        used_indices = used_indices + rec_group_hist

        random_rec = self._get_random_places(
            list(self.embed_dict.keys()), used_indices, parts[3])
        used_indices = used_indices + random_rec

        final_recomendation = query_rec_list + \
            rec_user_hist + rec_group_hist + random_rec
        random.shuffle(final_recomendation)
        return final_recomendation

    def update_embed_dict(self):
        self.embed_dict = self._load_embeddings()
