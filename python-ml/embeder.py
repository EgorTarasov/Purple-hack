from langchain_community.embeddings import HuggingFaceEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from typing import Any
from pymongo import MongoClient


class Embeder:
    def __init__(self, constr: str, model="intfloat/multilingual-e5-large") -> None:

        self.mongo = MongoClient(constr)

        db = self.mongo["cbr"]
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=512, chunk_overlap=0
        )
        splits = self.text_splitter.create_documents(data, metadatas=metadata)
        self.materials = db["materials"]
        self.embed_model = HuggingFaceEmbeddings(model_name=model)

    def need_embedings(self, doc_id: int, page: int) -> bool:
        """return true if document has embedings or does not exists"""
        res = self.materials.find_one({"docID": str(doc_id), "pageID": str(page)})
        if res is None:
            return False
        elif "embedings" in res.keys():
            return False
        else:
            return True

    def get_embeds(self, txt: str) -> list[Any]:

        return list()
