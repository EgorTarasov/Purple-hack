import sys
import os
import requests
import json
from typing import Any
import logging
from concurrent import futures
from typing import Generator
import search_engine_pb2_grpc
from search_engine_pb2 import Query, Response
from grpc import ServicerContext, server

from sentence_transformers import SentenceTransformer

from store import VectorStore
from llm import Ollama

log = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)


class MLservice:

    def __init__(
        self,
        clickhouse_uri: str,
        ollama_uri: str,
        embedding_model: str = "Tochka-AI/ruRoPEBert-e5-base-2k",
        table_name: str = "vector_search",
    ) -> None:
        model = SentenceTransformer(embedding_model)

        self.store = VectorStore(clickhouse_uri, model, table_name=table_name)
        self.ollama = Ollama(ollama_uri)

    def get_response(self, query: str) -> tuple[str, str]:
        """return response and context from ollama model"""
        docs = self.store.search_similarity(query, 3)
        return self.ollama.get_response(query, docs)

    def get_stream_response(self, query: str) -> Generator[str, Any, None]:
        """return response and context from ollama model"""
        docs = self.store.search_similarity(query, 3)
        for chunk in self.ollama.get_stream_response(query, docs):
            yield chunk["message"]["content"]  # type: ignore
        yield "Ссылки:" + "\n".join([doc.metadata for doc in docs])


class SearchEngingeServicer(search_engine_pb2_grpc.SearchEngineServicer):
    def __init__(self, service: MLservice) -> None:
        self.service = service
        super().__init__()

    def Respond(self, query: Query, ctx: ServicerContext) -> Response:
        log.debug(f"query = {query}")
        log.debug(f"ctx = {ctx}")
        body, context = self.service.get_response(query.body)

        return Response(body=body, context=context)

    def RespondStream(
        self, query: Query, ctx: ServicerContext
    ) -> Generator[Response, Any, None]:
        body: str = query.body
        log.info(f"new request: {body}")
        model: str = query.model
        for chunk in self.service.get_stream_response(body):
            if chunk.startswith("Ссылки:"):
                yield Response(body="", context=chunk)
            else:
                yield Response(body=chunk, context="")


def serve():
    s = server(futures.ThreadPoolExecutor(max_workers=10))
    clickhouse_uri = (
        "clickhouse://testuser:superstrongpassword@larek.tech:65002/default"
    )
    ollama_uri = "http://larek.tech:11434"
    if clickhouse_uri is None:
        log.error("CLICKHOUSE_URI is not set")
        sys.exit(1)
    if ollama_uri is None:
        log.error("OLLAMA_URI is not set")
        sys.exit(1)

    ml = MLservice(
        clickhouse_uri=clickhouse_uri,
        ollama_uri=ollama_uri,
    )

    search_engine_pb2_grpc.add_SearchEngineServicer_to_server(
        SearchEngingeServicer(ml), s
    )
    log.info("starting server")
    s.add_insecure_port("[::]:10000")
    s.start()
    s.wait_for_termination()


if __name__ == "__main__":
    serve()
