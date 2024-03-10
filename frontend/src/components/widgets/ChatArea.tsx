import { useContext, useEffect, useRef, useState } from "react";
import { ScrollArea } from "../ui/scroll-area";
import { Textarea } from "../ui/textarea";
import { Button } from "../ui/button";
import Message from "../entities/Message";
import { WebsocketContext, WebsocketContextType } from "@/context/WebSocketProvider";
import { IMessage } from "@/models";
import uuid from "react-uuid";

function ChatArea() {
	const [messageList, setMessageList, ready, val, send] : WebsocketContextType = useContext(WebsocketContext);

	const [currentMessage, setCurrentMessage] = useState("");

	const sendMessage = async () => {
		if (currentMessage !== "") {
			const messageData: IMessage = {
				data: currentMessage,
				time:
					new Date(Date.now()).getHours() +
					":" +
					new Date(Date.now()).getMinutes(),
				senderChat: false,
				id: uuid(),
			};

			if (ready) {
				send("message");
			}

			setMessageList((list) => [...list, messageData]);
			setCurrentMessage("");
		}
	};

	useEffect(() => {
		const messageData: IMessage = {
			data: val,
			time:
				new Date(Date.now()).getHours() +
				":" +
				new Date(Date.now()).getMinutes(),
			senderChat: true,
			id: uuid(),
		};
		setMessageList((list) => [...list, messageData]);
	}, [val, setMessageList]);

	const scrollRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		if (scrollRef.current) {
			const scrollAreaNode = scrollRef.current;
			scrollAreaNode.scrollTop = scrollAreaNode.scrollHeight;
		}
	}, [messageList]);

	return (
		<div className="w-full px-10 py-4">
			<ScrollArea
				className="h-[calc(100vh-300px)] p-5 mb-5 rounded-md border bg-white"
				scrollRef={scrollRef}
			>
				{messageList.map((messageContent) => {
					return <Message key={messageContent.id} message={messageContent} />;
				})}
			</ScrollArea>
			<form className="mx-20 max-h-[200px]">
				<div className="grid gap-4">
					<Textarea
						value={currentMessage}
						className="p-4"
						placeholder={`Ваш вопрос ...`}
						onChange={(event) => {
							setCurrentMessage(event.target.value);
						}}
						onKeyDown={(event) => {
							event.key === "Enter" && sendMessage();
						}}
					/>
					<div className="flex items-center">
						<Button type="button" onClick={sendMessage} size="sm" className="ml-auto">
							Отправить
						</Button>
					</div>
				</div>
			</form>
		</div>
	);
}

export default ChatArea;
