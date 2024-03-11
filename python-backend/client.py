import grpc
from search_engine_pb2_grpc import SearchEngineStub
from search_engine_pb2 import Query


def run():
    # Open a gRPC channel
    with grpc.insecure_channel("localhost:10000") as channel:

        # Create a stub (client)
        stub = SearchEngineStub(channel)

        # Create a valid request

        query = Query(
            body="Каковы критерии определения достаточности денежных средств на банковском счете плательщика, учитывая остаток средств на начало текущего дня и суммы, необходимые для учета?",
            model="Your model name",
        )

        # Make the call
        # response = stub.Respond(query)

        # # Print the response
        # print(response)

        responses = stub.RespondStream(query)

        # # Iterate over the stream of responses
        for response in responses:
            # Print the response
            print(response)


if __name__ == "__main__":
    run()