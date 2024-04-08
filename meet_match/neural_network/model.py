import torch
import pandas as pd
import torch.nn.functional as F

from torch import Tensor
from transformers import AutoTokenizer, AutoModel


tokenizer = AutoTokenizer.from_pretrained("intfloat/multilingual-e5-large-instruct")
model = AutoModel.from_pretrained("intfloat/multilingual-e5-large-instruct")


def average_pool(last_hidden_states: Tensor, attention_mask: Tensor) -> Tensor:
    last_hidden = last_hidden_states.masked_fill(~attention_mask[..., None].bool(), 0.0)
    return last_hidden.sum(dim=1) / attention_mask.sum(dim=1)[..., None]


def get_detailed_instruct(task_description: str, query: str) -> str:
    return f"Instruct: {task_description}\nQuery: {query}"


def search_places(
    queries: list[str],
    passages: list[torch.Tensor],
    from_place: int,
    number_of_places: int,
) -> list[list[tuple[int, float]]]:
    task = "Represent this sentence for searching relevant passages: "
    queries_detailed = [get_detailed_instruct(task, query) for query in queries]

    queries_tokenized = tokenizer(
        queries_detailed,
        max_length=512,
        padding=True,
        truncation=True,
        return_tensors="pt",
    )
    with torch.no_grad():
        outputs = model(**queries_tokenized)

    embeddings = average_pool(
        outputs.last_hidden_state, queries_tokenized["attention_mask"]
    )
    query_embeddings = F.normalize(embeddings, p=2, dim=1)
    scores = (query_embeddings @ passages.T) * 100

    sorted_scores, sorted_indices = torch.sort(scores, descending=True)

    passages_for_queries = []
    for i in range(len(queries)):

        top_scores = sorted_scores[i][from_place : from_place + number_of_places]
        top_indices = sorted_indices[i][from_place : from_place + number_of_places]
        passages_for_queries.append(
            [
                (index.item(), score.item())
                for index, score in zip(top_indices, top_scores)
            ]
        )

    return passages_for_queries


def predict(csv_file, text_embeddings, query, from_line, to_line):
    list_res = search_places(
        query.split(), text_embeddings, from_line, to_line - from_line + 1
    )

    return list_res
