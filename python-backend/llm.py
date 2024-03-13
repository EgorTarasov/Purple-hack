from typing import Any
from typing import Iterator, Mapping
from store import result_document
import json
from ollama import Client


class Ollama:
    def __init__(
        self,
        uri: str = "http://localhost:11434/api/generate",
        template: str | None = None,
    ) -> None:
        self.client = Client(host=uri)
        if template is None:
            self.template = """Отвечай только на русском. Если пишешь на другом языке, переводи его на русской.
Если не знаешь ответа, скажи что не знаешь ответа, не пробуй отвечать.
Я дам тебе три текста, из которых надо дать ответ на поставленный вопрос.
Также тебе надо оставить ссылку из источник.

Context:
источник {url1}:
{context1}

источник {url2}:
{context2}

источник {url3}:
{context3}

Вопрос: {question} на русском языке. Ответь на вопрос основываясь на данных документах
Развернутый ответ:
"""
        else:
            self.template = template

    def _get_prompt(self, query: str, docs: list[result_document]) -> str:
        return self.template.format(
            context1=docs[0].text,
            url1=docs[0].metadata,
            context2=docs[1].text,
            url2=docs[1].metadata,
            context3=docs[2].text,
            url3=docs[2].metadata,
            question=query,
        )

    def get_response(self, query: str, docs: list[result_document]) -> tuple[str, str]:

        prompt = self._get_prompt(query, docs)
        print(prompt)
        response = self.client.chat(
            model="llama2",
            messages=[{"role": "user", "content": prompt}],
            options={"temperature": 0},
        )

        return response["message"]["content"], "Ссылки:" + "\n".join([doc.metadata for doc in docs])  # type: ignore

    def get_stream_response(
        self, query: str, docs: list[result_document]
    ) -> Mapping[str, Any] | Iterator[Mapping[str, Any]]:
        prompt = self._get_prompt(query, docs)
        stream = self.client.chat(
            model="llama2",
            messages=[{"role": "user", "content": prompt}],
            options={"temperature": 0},
            stream=True,
        )
        return stream
