import { useEffect, useRef, useState } from "react";
import { ScrollArea } from "../ui/scroll-area";
import { Textarea } from "../ui/textarea";
import { Button } from "../ui/button";
import Message from "../entities/Message";
import { WebsocketContextType, useWS } from "@/context/WebSocketProvider";
import { IMessage } from "@/models";
import uuid from "react-uuid";
import { streamIndicator } from "../../constants";

function ChatArea() {
	const maxLengthSymbols = 5000;
	const [lengthSymbols, setLengthSymbols] = useState(0);

	const [messageList, setMessageList, ready, val, send]: WebsocketContextType =
		useWS();

	const [currentMessage, setCurrentMessage] = useState("");
	const [isStreamStarted, setIsStreamStarted] = useState(false);

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
				error: false,
			};

			if (ready) {
				send(currentMessage);
			}

			setMessageList((list) => [...list, messageData]);
			setCurrentMessage("");
		}
	};

	const sendCancel = async () => {
		send(streamIndicator.cancel)
		// setIsStreamStarted(false);
	};

	useEffect(() => {
		function mutateMessageList(index: number, val: string, error: boolean) {
			const newMessageList: IMessage[] = messageList.map((currentMessage, i) => {
				if (i === index) {
					const newMessage = currentMessage;
					newMessage.data += val;
					newMessage.error = error;
					return newMessage
				} else {
					return currentMessage;
				}
			});
			setMessageList(newMessageList);
		}

		let messageData: IMessage = {
			data: "",
			time: convertTime(),
			senderChat: true,
			id: uuid(),
			error: false,
		};


		if (!isStreamStarted && val !== "" && val !== streamIndicator.error && val !== streamIndicator.finished) {
			messageData = {
				data: val,
				time: convertTime(),
				senderChat: true,
				id: uuid(),
				error: false,
			};
			setMessageList((list) => [...list, messageData]);
			setIsStreamStarted(true);
		}
		else if (isStreamStarted && val !== streamIndicator.error && val !== streamIndicator.finished) {
			mutateMessageList(messageList.length - 1, val, false);
		}
		else if(isStreamStarted && val === streamIndicator.error){
			setIsStreamStarted(false);
			mutateMessageList(messageList.length - 1, "", true);
		}
		else if(isStreamStarted && val === streamIndicator.finished){
			setIsStreamStarted(false);
			mutateMessageList(messageList.length - 1, "", false);
		}
		else if(!isStreamStarted && val === streamIndicator.error){
			messageData = {
				data: "",
				time: convertTime(),
				senderChat: true,
				id: uuid(),
				error: true,
			};
			setMessageList((list) => [...list, messageData]);
		}
	}, [val]);

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
		<div className="mx-10 my-4">
			<ScrollArea
				className="h-[calc(100vh-250px)] p-5 mb-5 border-none bg-white"
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
						disabled={isStreamStarted}
					/>
					<div className="flex items-center">
						<p>
							{lengthSymbols}/{maxLengthSymbols}
						</p>
						{!isStreamStarted && <Button
							type="button"
							onClick={sendMessage}
							size="sm"
							className="ml-auto"
						>
							Отправить
						</Button>}
						{isStreamStarted && <Button
							type="button"
							onClick={sendCancel}
							size="sm"
							className="ml-auto"
						>
							Отменить генерацию ответа
						</Button>}
					</div>
				</div>
			</form>
		</div>
	);
}

export default ChatArea;
