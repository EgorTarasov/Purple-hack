import SideBar from "@/components/widgets/SideBar";
import ChatArea from "@/components/widgets/ChatArea";
import { WebsocketProvider } from "@/context/WebSocketProvider";
import { useParams } from "react-router-dom";
import ModelSelect from "@/components/widgets/ModelSelect";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Download } from "lucide-react";
import { toast } from "@/components/ui/use-toast";

const Chat = () => {
	const { id } = useParams<{ id: string }>();
	const [selectedModel, setSelectedModel] = useState<string>("llama");

	useEffect(() => {}, [selectedModel]);

	const handleModelSelectChange = (value: string) => {
		setSelectedModel(value);
	};

	function handleShareSessionButton(){
		const sharedLink = window.location.href.replace("chat", "history-shared");
		navigator.clipboard.writeText(sharedLink);
		toast({
			title: "Ссылка скопирована",
			description: (
				<div>
					<p>Вы можете поделиться историей чата</p>
					<p>
						<a href={sharedLink} target="_blank" style={{userSelect: 'all', textDecoration: 'underline'}}> 
							{sharedLink}
						</a>
					</p>
				</div>
			)
		});
		
	}

	return (
		<>
			<div className="h-full flex justify-between bg-white">
				{id && (
					<WebsocketProvider
						socketUuid={id}
						messageListDefault={[]}
						modelType={selectedModel}
					>
						<SideBar />
						<div className="grow">
							<div className="m-2 flex justify-between">
								<ModelSelect onSelectChange={handleModelSelectChange} />
								<Button onClick={handleShareSessionButton}>
									<Download />
								</Button>
							</div>
							<ChatArea />
						</div>
					</WebsocketProvider>
				)}
			</div>
		</>
	);
};

export default Chat;