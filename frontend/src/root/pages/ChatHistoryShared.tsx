import { useParams } from "react-router-dom";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useEffect, useState } from "react";
import { IMessage } from "@/models";
import Message from "@/components/entities/Message";
import ApiSession from "@/services/apiSession";
import { sessionToMessage } from "@/lib/sessionToMessage";

const ChatHistoryShared = () => {
	const { id } = useParams<{ id: string }>();

	const [messageSharedList, setMessageSharedList] = useState<IMessage[]>([]);

	useEffect(() => {
		async function fetchData() {
			try {
				// Получение сессий
				const sessions = await ApiSession.getUserSession();
				const historySessions = sessions.data;

				if (historySessions) {
					const newMessageLists: IMessage[][] = sessionToMessage(historySessions);
					setMessageSharedList(newMessageLists[0]);
				}
			} catch (error) {
				console.log(error);
			}
		}

		fetchData();
	}, []);

	return (
		<>
			<div className="h-full flex justify-between bg-white">
				{id && (
					<>
						<div className="grow">
							<div className="mx-10 my-4">
								Чат {id}
								<ScrollArea className="h-[calc(100vh-250px)] p-5 mb-5 border-none bg-white">
									{messageSharedList.map((messageContent) => {
										return (
											<Message
												key={messageContent.id}
												message={messageContent}
											/>
										);
									})}
								</ScrollArea>
							</div>
						</div>
					</>
				)}
			</div>
		</>
	);
};

export default ChatHistoryShared;
