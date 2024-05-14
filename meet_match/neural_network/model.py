"""
пример работы:

rec_sys = RecommendationSystem()

rec_sys.get_rec(user_id = 1, session_id = 1,"яблоко", num_idx = 10)
rec_sys.get_rec(user_id = 1,session_id = 1,"японская кухня", num_idx = 10, criteria = {"categories":['restaurants'])
rec_sys.get_rec(user_id = 1,session_id = 1,"японская кухня", num_idx = 10, criteria = {"categories":['restaurants'], "day":'Tuesday', 'time':'17:01'})


спустя некоторое время добавляем мвекторы новых мест в таблицу embeddings
rec_sys.build_embeddings()

и обновляем эмбеддинги и информацию о местах в самой модели

rec_sys.update_model_places()
"""

from datetime import datetime
from typing import List, Dict, Any, Tuple
from tqdm import tqdm_notebook
from transformers import AutoTokenizer, AutoModel
from torch import Tensor
import numpy as np
import ast
import json
import random
import psycopg2 as pg
import pickle
import torch.nn.functional as F
import pandas as pd
import torch
import logging
import logging.config
import logging.handlers

log = logging.getLogger(__name__)


def init_logging():
    """
    Инициализация логгера
    :return:
    """
    log_format = f"[%(asctime)s] [ Python server ] [%(levelname)s]:%(name)s:%(message)s"
    formatters = {"basic": {"format": log_format}}
    handlers = {"stdout": {
        "class": "logging.StreamHandler", "formatter": "basic"}}
    level = "INFO"
    handlers_names = ["stdout"]
    loggers = {
        "": {"level": level, "propagate": False, "handlers": handlers_names},
    }
    logging.basicConfig(level="INFO", format=log_format)
    log_config = {
        "version": 1,
        "disable_existing_loggers": False,
        "formatters": formatters,
        "handlers": handlers,
        "loggers": loggers,
    }
    logging.config.dictConfig(log_config)


init_logging()


class DatabaseManager:
    def __init__(self):
        self.connection_info = {
            "host": "proj_bd",
            "database": "meetmatch_db",
            "port": 5432,
            "user": "any1",
            "password": "1",
        }

    def _execute_query(self, query: str, params: Tuple = ()) -> List[Tuple]:
        try:
            with pg.connect(**self.connection_info) as conn:
                cursor = conn.cursor()
                cursor.execute(query, params)
                return cursor.fetchall()
        except pg.Error as e:
            print(f"An error occurred: {e}")
            return []

    def _execute_non_query(self, query: str, params: Tuple = ()):
        try:
            with pg.connect(**self.connection_info) as conn:
                cursor = conn.cursor()
                cursor.execute(query, params)
                conn.commit()
        except pg.Error as e:
            print(f"An error occurred: {e}")

    def get_swipes_for_session(
        self, user_id: int, session_id: int, for_group: bool = False
    ) -> List[Tuple[int, bool]]:
        if for_group:
            query = """
            SELECT place_id, is_liked
            FROM fact_scrolled
            WHERE user_id != %s AND session_id = %s
            """
        else:
            query = """
            SELECT place_id, is_liked
            FROM fact_scrolled
            WHERE user_id = %s AND session_id = %s
            """

        result = self._execute_query(query, (user_id, session_id))
        return [(int(place_id), bool(is_liked)) for place_id, is_liked in result]

    def get_all_descriptions(self):
        query = """
            SELECT place_id, title, description
            FROM places
        """
        result = self._execute_query(query)
        return [
            (int(place_id), str(title), str(description))
            for place_id, title, description in result
        ]

    def save_embedding(self, place_id, embedding):
        query = """
                    INSERT INTO embeddings (place_id, embedding) VALUES (%s,%s)
                """
        serialized_embedding = pickle.dumps(embedding)
        self._execute_non_query(query, (place_id, serialized_embedding))

    def load_embeddings(self) -> Dict[int, Any]:
        query = "SELECT place_id, embedding FROM embeddings"
        embeddings_data = self._execute_query(query)
        return {
            place_id: pickle.loads(embedding) for place_id, embedding in embeddings_data
        }

    def parse_timetable(self, timetable: str) -> Dict[str, str]:
        days_map = {
            'пн': 'Monday',
            'вт': 'Tuesday',
            'ср': 'Wednesday',
            'чт': 'Thursday',
            'пт': 'Friday',
            'сб': 'Saturday',
            'вс': 'Sunday',
            'ежедневно': 'Monday-Sunday'
        }
        timetable_dict = {}
        parts = timetable.split(',')
        for part in parts:
            try:
                days_part, *hours_part = part.strip().split(' ')
                days = days_part.split('–')
                hours = ' '.join(hours_part).strip()
                if len(days) == 2:
                    start_day, end_day = days_map[days[0]], days_map[days[1]]
                    current_day = start_day
                    while current_day != end_day:
                        timetable_dict[current_day] = hours
                        current_day = list(days_map.values())[
                            (list(days_map.values()).index(current_day) + 1) % 7]
                    timetable_dict[end_day] = hours
                elif days_part == 'ежедневно':
                    for day in days_map.values():
                        timetable_dict[day] = hours
                else:
                    timetable_dict[days_map[days[0]]] = hours
            except Exception as e:
                print(f"Error parsing timetable part '{part}': {e}")
                return {}

        return timetable_dict

    def fetch_places_info(self) -> Dict[int, Dict[str, Any]]:
        query = "SELECT place_id, categories, timetable FROM places"
        results = self._execute_query(query)
        places_info = {}
        for place_id, categories, working_hours in results:
            try:
                places_info[place_id] = {
                    "categories": ast.literal_eval(categories) if categories else [],
                    "working_hours": self.parse_timetable(working_hours) if working_hours else {}
                }
            except (ValueError, SyntaxError) as e:
                print(
                    f"Error parsing place info for place_id '{place_id}': {e}")
                return {}
        return places_info


class RecommendationSystem:
    def __init__(self):
        try:
            log.info("Init tokenizer from file")
            self.tokenizer = AutoTokenizer.from_pretrained(
                "./tokenizer_multil_e5_large"
            )
        except:
            log.info("Init tokenizer from net")
            self.tokenizer = AutoTokenizer.from_pretrained(
                "intfloat/multilingual-e5-large-instruct"
            )

        try:
            log.info("Init model from file")
            self.model = AutoModel.from_pretrained("./model_multil_e5_large")
        except:
            log.info("Init model from net")
            self.model = AutoModel.from_pretrained(
                "intfloat/multilingual-e5-large-instruct"
            )
        self.db_manager = DatabaseManager()
        self.embed_dict = self.db_manager.load_embeddings()
        self.places_info = self.db_manager.fetch_places_info()

    def _get_swipes_for_session(
        self, user_id: int, session_id: int, for_group: bool = False
    ) -> List[Tuple[int, bool]]:
        """Возвращает список пар (place_id, is_liked) для заданного пользователя и сессии, используя db_manager."""
        return self.db_manager.get_swipes_for_session(user_id, session_id, for_group)

    def load_embed_dict(self):
        """Обновляет словарь эмбеддингов, загружая данные через db_manager."""
        self.embed_dict = self.db_manager.load_embeddings()

    def load_places_info(self):
        """Обновляет информацию о местах, загружая данные через db_manager."""
        self.places_info = self.db_manager.fetch_places_info()

    def _is_open(self, working_hours, day, time_str):
        hours = working_hours.get(day, "")

        if hours == "fulltime":
            return True

        if not hours or "-" not in hours:
            return False

        open_time_str, close_time_str = hours.split("-")
        format_str = "%H:%M"  # Формат времени часы:минуты
        open_time = datetime.strptime(open_time_str, format_str).time()
        close_time = datetime.strptime(close_time_str, format_str).time()
        current_time = datetime.strptime(time_str, format_str).time()

        # Если время закрытия меньше времени открытия, возможно, заведение работает ночью.
        if close_time < open_time:
            return current_time >= open_time or current_time <= close_time
        else:
            return open_time <= current_time <= close_time

    def _filter_place(self, place_idx, criteria):

        # TODO: fix
        details = {}
        if self.places_info:
            details = self.places_info[place_idx]

        category_match = True
        if "categories" in criteria:
            category_match = any(
                cat in details["categories"] for cat in criteria["categories"]
            )

        time_match = True
        if "day" in criteria and "time" in criteria:
            time_match = self._is_open(
                details["working_hours"], criteria["day"], criteria["time"]
            )

        return category_match and time_match

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
        return f"Instruct: {task_description}\nQuery: {query}"

    def __average_pool(
        self, last_hidden_states: Tensor, attention_mask: Tensor
    ) -> Tensor:
        last_hidden = last_hidden_states.masked_fill(
            ~attention_mask[..., None].bool(), 0.0
        )
        return last_hidden.sum(dim=1) / attention_mask.sum(dim=1)[..., None]

    def _generate_embedding(self, text):
        inputs = self.tokenizer(
            text, return_tensors="pt", padding=True, truncation=True, max_length=512
        )
        with torch.no_grad():
            outputs = self.model(**inputs)

        # embedding = outputs.last_hidden_state.mean(dim=1).numpy()
        pooled_embeddings = self.__average_pool(
            outputs.last_hidden_state, inputs["attention_mask"]
        )
        embedding_normalized = F.normalize(pooled_embeddings, p=2, dim=1)
        return embedding_normalized

    def _get_rec_on_query(
        self, query: str, used_ids: List[int], criteria: Dict, num_items: int = 20
    ) -> List[int]:
        """
        Получает рекомендации на основе текстового запроса, исключая уже использованные ID.

        Параметры:
        - query (str): Текстовый запрос для поиска.
        - embeddings (Dict[int, torch.Tensor]): Словарь с эмбеддингами мест, где ключ - это ID места.
        - used_ids (List[int]): Список ID мест, которые уже были использованы и должны быть исключены из результатов.
        - criteria (Dict): словарь фильтров.
        - num_items (int): Количество возвращаемых рекомендаций.

        Возвращает:
        - List[int]: Список ID мест, которые рекомендуются на основе запроса.
        """

        task = "Represent this sentence for searching relevant passages: "
        query_embedding = self._generate_embedding(
            self.__get_detailed_instruct(task, query)
        )

        place_ids = list(self.embed_dict.keys())
        embedding_matrix = torch.cat(list(self.embed_dict.values()), dim=0)

        similarities = self._cosine_similarity(
            query_embedding, embedding_matrix
        ).squeeze(0)
        sorted_scores, sorted_indices = torch.sort(
            similarities, descending=True)

        top_n_place_ids = []
        idx_iter = 0
        while len(top_n_place_ids) < num_items and idx_iter < len(place_ids):
            current_idx = sorted_indices[idx_iter].item()

            if place_ids[current_idx] not in used_ids and self._filter_place(
                place_ids[current_idx], criteria
            ):
                top_n_place_ids.append(place_ids[current_idx])

            idx_iter += 1

        return top_n_place_ids

    def _generate_rec_on_user_session_hist(
        self,
        swiped_places: List[Tuple[int, bool]],
        used_ids: List[int],
        criteria,
        num_top_idx: int = 20,
        dislike_weight: float = 0.5,
    ) -> List[int]:
        """
        Генерирует рекомендации на основе истории лайков и дизлайков пользователя в данной сессии.

        Параметры:
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

        # Извлечение эмбеддингов для лайкнутых мест
        if liked_place_ids:
            liked_embeddings = torch.cat(
                [
                    self.embed_dict.get(i, torch.zeros((1, 1024)))
                    for i in liked_place_ids
                ],
                dim=0,
            )
        else:
            liked_embeddings = torch.zeros((1, 1024))

        # Извлечение эмбеддингов для дизлайкнутых мест
        if disliked_place_ids:
            disliked_embeddings = torch.cat(
                [
                    self.embed_dict.get(i, torch.zeros((1, 1024)))
                    for i in disliked_place_ids
                ],
                dim=0,
            )
        else:
            disliked_embeddings = torch.zeros((1, 1024))

        # Подготовка общей матрицы эмбеддингов и списка ID мест
        all_embeddings = torch.cat(list(self.embed_dict.values()), dim=0)
        place_ids = list(self.embed_dict.keys())

        # Расчёт косинусного сходства для лайкнутых и дизлайкнутых мест
        liked_similarities = self._cosine_similarity(
            liked_embeddings, all_embeddings)
        disliked_similarities = self._cosine_similarity(
            disliked_embeddings, all_embeddings
        )

        # Комбинирование сходства, учитывая вес дизлайков
        combined_scores = liked_similarities.mean(
            dim=0
        ) - dislike_weight * disliked_similarities.mean(dim=0)

        scores, indices = combined_scores.squeeze(0).sort(descending=True)

        top_n_place_ids = []
        idx_iter = 0
        while len(top_n_place_ids) < num_top_idx and idx_iter < len(indices):
            current_idx = indices[idx_iter].item()
            if place_ids[current_idx] not in used_ids and self._filter_place(
                place_ids[current_idx], criteria
            ):
                top_n_place_ids.append(place_ids[current_idx])
            idx_iter += 1

        return top_n_place_ids

    def _generate_group_rec_on_likes(
        self,
        user_id: int,
        session_id: int,
        swiped_places: List[Tuple[int, bool]],
        used_ids: List[int],
        criteria,
        num_top_idx: int = 20,
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
            idx for (idx, is_liked) in swiped_places if is_liked == 1
        ]

        if not group_liked_place_ids:
            return []
        group_liked_embeddings = torch.cat(
            [self.embed_dict[i] for i in group_liked_place_ids], dim=0
        )

        all_embeddings = torch.cat(list(self.embed_dict.values()), dim=0)
        place_ids = list(self.embed_dict.keys())

        group_liked_similarities = self._cosine_similarity(
            group_liked_embeddings, all_embeddings
        )

        combined_scores = group_liked_similarities.mean(dim=0)

        scores, indices = combined_scores.squeeze(0).sort(descending=True)

        top_n_place_ids = []
        idx_iter = 0
        while len(top_n_place_ids) < num_top_idx and idx_iter < len(indices):
            current_idx = indices[idx_iter].item()
            if place_ids[current_idx] not in used_ids and self._filter_place(
                place_ids[current_idx], criteria
            ):
                top_n_place_ids.append(place_ids[current_idx])
            idx_iter += 1

        return top_n_place_ids

    def _get_random_places(self, all_place_ids, used_indices, criteria, num_places=10):
        """
        Выбирает случайные места, не входящие в список использованных и соответствующие критериям фильтрации.

        Параметры:
        - all_place_ids (list): Список всех идентификаторов мест.
        - used_indices (list): Список уже использованных идентификаторов мест.
        - criteria (dict): Словарь критериев для фильтрации мест.
        - num_places (int): Количество мест для выборки.

        Возвращает:
        - list: Список идентификаторов мест, которые удовлетворяют критериям.
        """
        # Перемешиваем список индексов
        shuffled_place_ids = all_place_ids[:]
        random.shuffle(shuffled_place_ids)

        place_ids = []
        for idx in shuffled_place_ids:
            if len(place_ids) >= num_places:
                break
            if idx not in used_indices and self._filter_place(idx, criteria):
                place_ids.append(idx)

        return place_ids

    # сохраняет распределение вероятностей

    def _plan_recommendations(self, num_idx, rec_sources):
        probabilities = np.array(list(rec_sources.values()))

        probabilities /= probabilities.sum()

        counts = np.random.multinomial(num_idx, probabilities)

        return dict(zip(rec_sources.keys(), counts))

    def get_rec(
        self,
        user_id: int,
        session_id: int,
        query: str,
        num_idx=50,
        criteria: dict = None,
    ) -> List[int]:
        """
        Параметры:
        - user_id (int): Идентификатор пользователя.
        - session_id (int): Идентификатор сессии.
        - query (str): Текстовый запрос для поиска.
        - criteria (dict): Фильтры к запросу
        - num_idx (int): Число индексов на выходе.

        Возвращает:
        - List[int]: Список ID рекомендованных мест.
        """
        if not criteria:
            criteria = {}

        swipe_user_hist = self._get_swipes_for_session(
            user_id, session_id, for_group=False
        )

        swiped_places_count = len(swipe_user_hist)

        rec_sources = {
            'query_based': 0.4,  # поначалу только те, что по описанию
            'user_history_based': 0.2 if swiped_places_count >= 10 else 0.0,
            'group_based': 0.3 if swiped_places_count >= 10 else 0.0,
            'random': 0.05 if swiped_places_count >= 20 else 0.0
        }
        rec_plan = self._plan_recommendations(num_idx, rec_sources)

        used_indices = swipe_user_hist.copy()
        used_indices = [i[0] for i in used_indices]

        query_rec_list = self._get_rec_on_query(
            query, used_indices, criteria, rec_plan["query_based"]
        )
        used_indices = used_indices + query_rec_list

        rec_user_hist = self._generate_rec_on_user_session_hist(
            swipe_user_hist,
            used_indices,
            criteria,
            rec_plan["user_history_based"],
        )
        used_indices = used_indices + rec_user_hist

        swipe_group_hist = self._get_swipes_for_session(
            user_id, session_id, for_group=True
        )

        rec_group_hist = self._generate_group_rec_on_likes(
            user_id,
            session_id,
            swipe_group_hist,
            used_indices,
            criteria,
            rec_plan["group_based"],
        )
        used_indices = used_indices + rec_group_hist

        random_rec = self._get_random_places(
            list(self.embed_dict.keys()
                 ), used_indices, criteria, rec_plan["random"]
        )
        used_indices = used_indices + random_rec

        final_recomendation = query_rec_list + rec_group_hist + random_rec

        if len(final_recomendation) < num_idx:
            additional_needed = num_idx - len(final_recomendation)
            additional_recommendations = self._get_rec_on_query(
                query, used_indices, criteria, additional_needed
            )
            final_recomendation.extend(additional_recommendations)

        random.shuffle(final_recomendation)

        return final_recomendation

    def build_embeddings(self):
        places = self.db_manager.get_all_descriptions()

        for place_id, title, description in places:
            title = title if title else ""
            description = description if description else ""

            text = "Название места: " + title + "   Описание: " + description
            embedding = self._generate_embedding(text)
            #             serialized_embedding = pickle.dumps(embedding)

            # Insert the embedding into the database
            self.db_manager.save_embedding(place_id, embedding)

    def build_embeddings(self):
        places = self.db_manager.get_all_descriptions()

        for place_id, title, description in places:
            if place_id in self.embed_dict:
                continue

            title = title if title else ""
            description = description if description else ""

            text = "Название места: " + title + "   Описание: " + description
            embedding = self._generate_embedding(text)

            # Insert the embedding into the database
            self.db_manager.save_embedding(place_id, embedding)

    def update_model_places(self):
        self.load_embed_dict()
        self.load_places_info()
