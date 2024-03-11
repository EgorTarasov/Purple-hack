from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Query(_message.Message):
    __slots__ = ("body", "model")
    BODY_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    body: str
    model: str
    def __init__(self, body: _Optional[str] = ..., model: _Optional[str] = ...) -> None: ...

class Response(_message.Message):
    __slots__ = ("body", "context")
    BODY_FIELD_NUMBER: _ClassVar[int]
    CONTEXT_FIELD_NUMBER: _ClassVar[int]
    body: str
    context: str
    def __init__(self, body: _Optional[str] = ..., context: _Optional[str] = ...) -> None: ...
