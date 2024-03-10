import SideBar from "@/components/widgets/SideBar";
import ChatArea from "@/components/widgets/ChatArea";


const Chat = () => {
	return (
		<>
			<div className="flex justify-normal">
				<SideBar />
				<ChatArea />
			</div>
		</>
	);
};

export default Chat;
