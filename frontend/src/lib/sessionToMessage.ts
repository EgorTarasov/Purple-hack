import { IMessage, IQuery, IResponses, ISession } from "@/models";

export function sessionToMessage(userSessions: ISession[]): IMessage[][] {
	const newMessageLists: IMessage[][] = [];
	userSessions.map((session) => {
		const newMessageList: IMessage[] = session.queries.flatMap(
			(query: IQuery, index: number) => {
				let response: IResponses;
				if (session.responses[index]) {
					response = session.responses[index]
					return [
						{
							id: query.id.toString(),
							senderChat: false,
							data: query.body,
							time: query.createdAt,
							error: false,
						},
						{
							id: response.id.toString(),
							senderChat: true,
							data: response.body,
							time: response.createdAt,
							error: false,
						},
					];
				}
				else {
					return [
						{
							id: query.id.toString(),
							senderChat: false,
							data: query.body,
							time: query.createdAt,
							error: false,
						},
					];
				}
			}
		);
		newMessageLists.push(newMessageList);
	});
	return newMessageLists;
}
