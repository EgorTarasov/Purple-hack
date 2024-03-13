import { ScrollArea } from "../ui/scroll-area";
import Message from "../entities/Message";
import { useAuth } from "@/context/Authprovider";

function ChatAreaHistory() {
	const { messageHistoryListCurrent } = useAuth();

	return (
		<div className="mx-10 my-4">
			<ScrollArea
				className="h-[calc(100vh-250px)] p-5 mb-5 border-none bg-white"
			>
				{messageHistoryListCurrent.map((messageContent) => {
					return <Message key={messageContent.id} message={messageContent} />;
				})}
			</ScrollArea>
		</div>
	);
}

export default ChatAreaHistory;
