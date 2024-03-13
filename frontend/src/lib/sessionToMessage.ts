import { IMessage, IQuery, IResponses, ISession } from "@/models";

export function sessionToMessage(userSessions: ISession[]) : IMessage[][] {
	const newMessageLists: IMessage[][] = [];
	userSessions.map((session) => {
		const newMessageList: IMessage[] = session.queries.flatMap(
			(query: IQuery, index: number) => {
				const response: IResponses = session.responses[index];
				return [
					{
						id: query.id.toString(),
						senderChat: false,
						data: query.body,
						time: query.createdAt.toISOString(),
						error: false,
					},
					{
						id: response.id.toString(),
						senderChat: true,
						data: response.body,
						time: response.createdAt.toISOString(),
						error: false,
					},
				];
			}
		);
		newMessageLists.push(newMessageList);
	});
    return newMessageLists;
}
