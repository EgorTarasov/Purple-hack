import { IMessage } from "@/models";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import logoImg from "../../assets/OnlyLogo.svg";

interface IMessageProps {
	message: IMessage;
}
const Message = ({ message }: IMessageProps) => {
	return (
		<div className="my-5">
			<div className="flex items-center gap-2">
				<Avatar>
					<AvatarImage src={message.senderChat? logoImg : "https://github.com/shadcn.png"} alt="@shadcn" />
					<AvatarFallback>{message.senderChat ? "Чат" : "Вы"}</AvatarFallback>
				</Avatar>
				<p className="font-bold">{message.senderChat ? "Чат" : "Вы"}</p>
			</div>
			<div>{message.data}</div>
		</div>
	);
};

export default Message;
