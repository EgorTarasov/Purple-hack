import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "../ui/button";
import { useParams, useNavigate, Link } from "react-router-dom";
import uuid from "react-uuid";
import { useToast } from "../ui/use-toast";
import { ToastAction } from "../ui/toast";
import { WebsocketContextType, useWS } from "@/context/WebSocketProvider";
import logoPath from "../../assets/MediumLogo.svg";
import CreateAccountDialog from "./CreateAccountDialog";
import { useAuth } from "@/context/Authprovider";
import { IMessage } from "@/models";

export default function SideBar() {
	const { id } = useParams<{ id: string }>();

	const navigate = useNavigate();
	const { toast } = useToast();
	const { isAuthorized, messageHistoryLists, setMessageHistoryListCurrent } =
		useAuth();

	const [messageList]: WebsocketContextType = useWS();

	function handleNewChat() {
		if (!isAuthorized) {
			toast({
				title: "История чата не сохранится",
				description: "Для сохранения истории необходимо войти",
				action: (
					<ToastAction
						altText="Новый чат"
						onClick={() => {
							navigate(`/chat/${uuid()}`);
						}}
					>
						Все равно создать
					</ToastAction>
				),
			});
		}
	}

	function handleHistoryChatButton(messageListParam: IMessage[]) {
		setMessageHistoryListCurrent(messageListParam);
		navigate(`/history/${uuid()}`);
	}

	return (
		<div className="border-r p-4 border-border-color w-[320px] bg-gradient-to-br from-secondary-medium to-white flex flex-col">
			<div className="flex gap-3 items-center mb-2">
				<Link to="/">
					<img src={logoPath} alt="logo" width={130} height={325} />
				</Link>
			</div>
			<Button className="w-[100%] mb-2 bg-border-color" onClick={handleNewChat}>
				Новый чат
			</Button>
			<ScrollArea className="w-100% grow">
				<div className="my-4">Активный чат</div>
				<Button variant="secondary" className="w-full p-2 mb-2" key={id}>
					{messageList[0] ? messageList[0].data.substring(0, 35) : id}
				</Button>
				{isAuthorized && (
					<>
						<div className="my-4">История чатов</div>
						{messageHistoryLists.map((messageListParam) => {
							return (
								<Button
									variant="secondary"
									className="w-full p-2 mb-2"
									key={id}
									onClick={() => handleHistoryChatButton(messageListParam)}
								>
									{messageListParam[0].data.substring(0, 35)}
								</Button>
							);
						})}
					</>
				)}
			</ScrollArea>
			<CreateAccountDialog />
		</div>
	);
}
