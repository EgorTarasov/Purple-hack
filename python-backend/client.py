import grpc
import search_engine_pb2
import search_engine_pb2_grpc
from pprint import pprint


def run():
    # Create a channel to the server
    channel = grpc.insecure_channel("localhost:10000")

    # Create a stub (client)
    stub = search_engine_pb2_grpc.SearchEngineStub(channel)

    # Create a valid request message
    request = search_engine_pb2.Query(
        body="Какое отношение к денежно-кредитной политике имеет управление ликвидностью банковского сектора и ставками денежного рынка, которое осуществляет Банк России?",
        model="test model",
    )

    # Make the call to the Respond method
    # response = stub.Respond(request)

    # # Print the response
    # pprint(response.body)
    # pprint(response.context)

    # Make the call to the RespondStream method
    responses = stub.RespondStream(request)

    # Print the responses
    for response in responses:
        print(response)


if __name__ == "__main__":
    run()
