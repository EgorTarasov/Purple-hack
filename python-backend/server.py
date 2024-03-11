from pymongo import MongoClient
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_community.embeddings import HuggingFaceEmbeddings
from langchain.vectorstores.pgvector import PGVector
from langchain import PromptTemplate
from langchain.llms import Ollama
from langchain.chains import RetrievalQA

# for streaming
from langchain_core.runnables import RunnableParallel, RunnablePassthrough
from langchain_core.output_parsers import StrOutputParser
from langchain_community.chat_models import ChatOllama

import sys
import os

import logging
from concurrent import futures

import search_engine_pb2_grpc
from search_engine_pb2 import Query, Response
from grpc import ServicerContext, server


class DataLoader:

    def __init__(
            self,
            mongo_con: str,
            pg_con: str,
            collection_name: str = "ml_docs",
    ) -> None:

        self.mongo = MongoClient(mongo_con)
        db = self.mongo["cbr"]
        self.pg_con = pg_con
        self.collection_name = collection_name
        self.materials = db["materials"]
        self.embeddings = HuggingFaceEmbeddings(
            model_name="intfloat/multilingual-e5-large"
        )

    def get_store(self):  # -> Any:
        return PGVector(
            connection_string=self.pg_con,
            embedding_function=self.embeddings,
            collection_name="ml_docs",
        )

    def get_embedding(self, text: str):
        db = self.mongo["cbr"]

        materials = db["materials"]
        res = materials.find()
        data = [obj["text"] for obj in res]
        metadata = [{"link": f"[источник]({obj['src']})"} for obj in res]

        text_splitter = RecursiveCharacterTextSplitter(chunk_size=512, chunk_overlap=0)
        docs = text_splitter.create_documents(data, metadatas=metadata)
        embeddings = HuggingFaceEmbeddings(model_name="intfloat/multilingual-e5-large")
        db = PGVector.from_documents(
            embedding=embeddings,
            documents=docs,
            collection_name=self.collection_name,
            connection_string=self.pg_con,
        )


class SuppressStdout:
    def __enter__(self):
        self._original_stdout = sys.stdout
        self._original_stderr = sys.stderr
        sys.stdout = open(os.devnull, "w")
        sys.stderr = open(os.devnull, "w")

    def __exit__(self, exc_type, exc_val, exc_tb):
        sys.stdout.close()
        sys.stdout = self._original_stdout
        sys.stderr = self._original_stderr


class SearchEngingeServicer(search_engine_pb2_grpc.SearchEngineServicer):
    def __init__(self, loader: DataLoader, ollama_api: str) -> None:
        # TODO: add models params to params
        vectorstore = loader.get_store()
        self.ollama = Ollama(base_url=ollama_api, model="llama2", temperature=0)
        self.chatOllama = ChatOllama(base_url=ollama_api, model="llama2", temperature=0)

        template = """Отвечай только на русском. Если пишешь на другом языке, переводи его на русской.
Если не знаешь ответа, скажи что не знаешь ответа, не пробуй отвечать.
Я дам тебе пять текстов, из которых надо дать ответ на поставленный вопрос.

Context:
{context}

Question: {question} на русском языке.
Ответ:
"""
        prompt = PromptTemplate.from_template(
            template=template,
        )

        retriever = vectorstore.as_retriever(
            search_type="similarity", search_kwargs={"k": 5}
        )
        self.qa_chain = RetrievalQA.from_chain_type(
            self.ollama,
            return_source_documents=True,
            retriever=retriever,
            chain_type_kwargs={"prompt": prompt},
        )

        def format_docs(docs):
            return "\n\n".join(doc.page_content for doc in docs)

        rag_chain_from_docs = (
            # RunnablePassthrough.assign(context=(lambda x: format_docs(x["context"])))
                prompt
                | self.chatOllama
                | StrOutputParser()
        )
        self.rag_chain_with_source = RunnableParallel(
            {
                "context": retriever,
                "question": RunnablePassthrough(),
            }
        ).assign(answer=rag_chain_from_docs)
        super().__init__()

    def Respond(self, query: Query, ctx: ServicerContext):
        log.debug(f"query = {query}")
        log.debug(f"ctx = {ctx}")
        result = self.qa_chain({"query": query.body})
        # call to llama
        log.info(f"result = {result}")
        response_body = result["result"]
        return Response(body=response_body, context="test context")

    def RespondStream(self, query: Query, ctx: ServicerContext):
        body: str = query.body
        model: str = query.model

        for chunk in self.rag_chain_with_source.stream(body):
            if "question" in chunk:
                continue
            if "context" in chunk:
                logging.info(f"context {chunk['context']}")
                # context = "\n".join([objd])
            else:
                body = chunk["answer"]
                logging.info(chunk["answer"])

            yield Response(body=body, context="456")


def serve():
    s = server(futures.ThreadPoolExecutor(max_workers=10))
    mongo_con = os.getenv("MONGO_CONN_STR")
    pg_con = os.getenv("PG_CONN_STR")
    ollama_api = os.getenv("OLLAMA_API")
    if mongo_con is None or pg_con is None or ollama_api is None:
        exit(-1)

    l = DataLoader(mongo_con=mongo_con, pg_con=pg_con)
    search_engine_pb2_grpc.add_SearchEngineServicer_to_server(
        SearchEngingeServicer(l, ollama_api), s
    )
    logging.info("created ml")
    s.add_insecure_port("[::]:10000")
    logging.info("starting server")
    s.start()
    s.wait_for_termination()


logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)

if __name__ == "__main__":
    serve()