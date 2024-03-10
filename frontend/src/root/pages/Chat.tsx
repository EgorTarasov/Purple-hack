import SideBar from "@/components/widgets/SideBar";
import ChatArea from "@/components/widgets/ChatArea";
import { WebsocketProvider } from "@/context/WebSocketProvider";
import { useParams } from "react-router-dom";
import ModelSelect from "@/components/widgets/ModelSelect";
import { useEffect, useState } from "react";

const Chat = () => {
	const { id } = useParams<{ id: string }>();
	const [selectedModel, setSelectedModel] = useState<string>("llama");

	useEffect(() => {
	}, [selectedModel]);

	const handleModelSelectChange = (value: string) => {
		setSelectedModel(value);
	};

	return (
		<>
			<div className="flex justify-normal">
				<SideBar />
				<div className="ml-2 pt-4">
					<ModelSelect onSelectChange={handleModelSelectChange} />
				</div>
				{id && (
					<WebsocketProvider socketUuid={id} messageListDefault={[]}>
						<ChatArea />
					</WebsocketProvider>
				)}
			</div>
		</>
	);
};

export default Chat;
