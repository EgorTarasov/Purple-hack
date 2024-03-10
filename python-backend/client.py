# Python
import grpc
import search_engine_pb2
import search_engine_pb2_grpc


def run():
    # Create a gRPC channel
    channel = grpc.insecure_channel("localhost:10000")

    # Create a stub (client)
    stub = search_engine_pb2_grpc.SearchEngineStub(channel)

    # Create a valid request message
    model = "my_model"
    query = "Каковы критерии определения достаточности денежных средств на банковском счете плательщика, учитывая остаток средств на начало текущего дня и суммы, необходимые для учета?"
    search_request = search_engine_pb2.Query(body=query, model=model)

    # Make the call
    response = stub.Respond(search_request)

    print(f"Response: {response}")


if __name__ == "__main__":
    run()