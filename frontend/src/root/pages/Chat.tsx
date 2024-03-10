import SideBar from "@/components/widgets/SideBar";
import ChatArea from "@/components/widgets/ChatArea";
import { WebsocketProvider } from "@/context/WebSocketProvider";
import { useParams } from "react-router-dom";

const Chat = () => {
	const { id } = useParams<{ id: string }>();
	return (
		<>
			<div className="flex justify-normal">
				<SideBar />
				{id && <WebsocketProvider socketUuid={id}>
					<ChatArea />
				</WebsocketProvider>}
			</div>
		</>
	);
};

export default Chat;
