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
import ApiSession from "@/services/apiSession";
import { sessionToMessage } from "@/lib/sessionToMessage";

export default function SideBar() {
	const { id } = useParams<{ id: string }>();

	const navigate = useNavigate();
	const { toast } = useToast();
	const { isAuthorized, messageHistoryLists, setMessageHistoryListCurrent, setUserSessions, setMessageHistoryLists} =
		useAuth();

	const [messageList]: WebsocketContextType = useWS();

	async function getSessions(){
		try {
			const sessions = await ApiSession.getUserSession();
			setUserSessions(sessions.data);
			console.log('newsessions', sessions.data)

			if (sessions.data) {
				const newMessageLists: IMessage[][] = sessionToMessage(sessions.data);
				setMessageHistoryLists(newMessageLists);
			}
		} catch (error) {
			console.log(error);
		}

	}

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
		else {
			getSessions();
			navigate(`/chat/${uuid()}`);
		}
	}

	function handleHistoryChatButton(messageListParam: IMessage[]) {
		setMessageHistoryListCurrent(messageListParam);
		navigate(`/history/${messageListParam[0].sessionUuid}`);
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
				{messageList.length !== 0 && (
					<>
						<div className="my-4">Активный чат</div>
						<Button variant="secondary" className="w-full p-2 mb-2" key={id}>
							{messageList[0] ? messageList[0].data.substring(0, 35) : id}
						</Button>
					</>
				)}
				{isAuthorized && (
					<>
						<div className="my-4">История чатов</div>
						{messageHistoryLists.map((messageListParam) => {
							return (
								<Button
									variant="secondary"
									className="w-full p-2 mb-2"
									key={messageListParam[0].id}
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
