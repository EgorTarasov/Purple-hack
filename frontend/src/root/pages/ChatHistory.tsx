import SideBar from "@/components/widgets/SideBar";
import { useParams } from "react-router-dom";
import ChatAreaHistory from "@/components/widgets/ChatAreaHistory";

const ChatHistory = () => {
	const { id } = useParams<{ id: string }>();

	return (
		<>
			<div className="h-full flex justify-between bg-white">
				{id && (
					<>
						<SideBar />
						<div className="grow">
							<ChatAreaHistory />
						</div>
					</>
				)}
			</div>
		</>
	);
};

export default ChatHistory;
