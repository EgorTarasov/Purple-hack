import { useEffect, useRef, useState } from "react";
import { ScrollArea } from "../ui/scroll-area";
import { Textarea } from "../ui/textarea";
import { Button } from "../ui/button";
import Message from "../entities/Message";
import {
	WebsocketContextType,
	useWS,
} from "@/context/WebSocketProvider";
import { IMessage } from "@/models";
import uuid from "react-uuid";

function ChatArea() {
	const maxLengthSymbols = 5000;
	const [lengthSymbols, setLengthSymbols] = useState(0);

	const [messageList, setMessageList, ready, val, send]: WebsocketContextType = useWS();

	const [currentMessage, setCurrentMessage] = useState("");

	function convertTime() {
		let hours = String(new Date(Date.now()).getHours());
		let minutes = String(new Date(Date.now()).getMinutes());
		if (hours.length === 1) hours = "0" + hours;
		if (minutes.length === 1) minutes = "0" + minutes;
		return hours + ":" + minutes;
	}

	const sendMessage = async () => {
		if (currentMessage.trim() !== "") {
			const messageData: IMessage = {
				data: currentMessage,
				time: convertTime(),
				senderChat: false,
				id: uuid(),
			};

			if (ready) {
				send(currentMessage);
			}

			setMessageList((list) => [...list, messageData]);
			setCurrentMessage("");
		}
	};

	useEffect(() => {
		const messageData: IMessage = {
			data: val,
			time: convertTime(),
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

	useEffect(() => {
		setLengthSymbols(currentMessage.length);
	}, [currentMessage]);

	return (
		<div className="grow mx-10 my-4">
			<ScrollArea
				className="h-[calc(100vh-300px)] p-5 mb-5 rounded-md border bg-white"
				scrollRef={scrollRef}
			>
				{messageList.map((messageContent) => {
					return <Message key={messageContent.id} message={messageContent} />;
				})}
			</ScrollArea>
			<form className="mx-auto max-h-[200px] md:w-[80%] lg:w-[60%] xl:w-[80%]">
				<div className="grid gap-4">
					<Textarea
						maxLength={maxLengthSymbols}
						value={currentMessage}
						className="p-4"
						placeholder={`Ваш вопрос ...`}
						onChange={(event) => {
							setCurrentMessage(event.target.value);
						}}
						onKeyDown={(event) => {
							if (event.key === "Enter" && event.ctrlKey) sendMessage();
						}}
					/>
					<div className="flex items-center">
						<p>
							{lengthSymbols}/{maxLengthSymbols}
						</p>
						<Button
							type="button"
							onClick={sendMessage}
							size="sm"
							className="ml-auto"
						>
							Отправить
						</Button>
					</div>
				</div>
			</form>
		</div>
	);
}

export default ChatArea;
