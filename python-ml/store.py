from clickhouse_driver import connect, Client
from sentence_transformers import SentenceTransformer
from tqdm import tqdm
from document import Document
from typing import NamedTuple
from collections import namedtuple
import pprint

result_document = namedtuple("result_document", ["text", "metadata", "embedding", "score"])

class VectorStore:
    
    def __init__(self, uri: str, embedding_model: SentenceTransformer, table_name: str = "vector_search") -> None:
        self.client = Client.from_url("clickhouse://testuser:superstrongpassword@larek.tech:65002/default")
        self.table_name = table_name
        self.model = embedding_model


    def create_table(self):
        self.client.execute("SET allow_experimental_annoy_index = 1")
        self.client.execute(f"""CREATE TABLE IF NOT EXISTS {self.table_name}
        (
        DocName String,
        Metadata String,
        embedding Array(Float32),
        INDEX ann_index_1 embedding TYPE annoy('cosineDistance')
        )
        ENGINE = MergeTree
        ORDER BY DocName;""")


    def create_embs(self, documents: list[Document]):
        embeddings = []

        for re in tqdm(documents):
            emb = self.model.encode(re.page_content)
            embeddings.append((re.page_content, re.metadata['src'], emb))

        self.client.execute(
            f'INSERT INTO {self.table_name} VALUES',
            embeddings)
        self.client.disconnect()


    def search_similarity(self,  query: str, k: int = 1) -> list[result_document]:
        emb_q = self.model.encode(query)

        res =   self.client.execute(f""" SELECT DocName, Metadata, embedding, L2Distance(embedding, %(emb)s) AS score
        FROM {self.table_name}
        ORDER BY score ASC
        LIMIT %(k)s
        """, {"emb": list(emb_q), "k": k})
        if res is None:
            return []
        
        return [result_document(text=row[0], metadata=row[1], embedding=row[2], score=row[3]) for row in res]
    
    
