import { IMessage } from "@/models";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import logoImg from "../../assets/OnlyLogo.svg";
import avaImg from "../../assets/avatar.jpg";

interface IMessageProps {
	message: IMessage;
}
const Message = ({ message }: IMessageProps) => {
	return (
		<div className="my-5 max-w-[100%]">
			<div className="flex">
				<Avatar>
					<AvatarImage
						src={message.senderChat ? logoImg : avaImg}
						alt="avatar"
					/>
					<AvatarFallback>{message.senderChat ? "Чат" : "Вы"}</AvatarFallback>
				</Avatar>
				<div className="w-full ml-2">
					<div className="flex mt-2">
						<p className="font-bold">{message.senderChat ? "Чат" : "Вы"}</p>
						<div className="ml-2 mt-0.7">{message.time}</div>
					</div>
					<pre style={{ wordWrap: "break-word", maxWidth: '100%', whiteSpace: 'pre-wrap' }}>{message.data}</pre>
				</div>
			</div>
		</div>
	);
};

export default Message;
