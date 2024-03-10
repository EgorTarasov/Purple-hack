import { IMessage } from "@/models";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import logoImg from "../../assets/OnlyLogo.svg";
import avaImg from "../../assets/avatar.jpg";

interface IMessageProps {
	message: IMessage;
}
const Message = ({ message }: IMessageProps) => {
	return (
		<div className="my-5 ">
			<div className="flex">
				<Avatar>
					<AvatarImage
						src={message.senderChat ? logoImg : avaImg}
						alt="@shadcn"
					/>
					<AvatarFallback>{message.senderChat ? "Чат" : "Вы"}</AvatarFallback>
				</Avatar>
				<div className="ml-2">
					<div className="flex mt-2">
						<p className="font-bold">
							{message.senderChat ? "Чат" : "Вы"}
						</p>
						<div className="ml-2 mt-0.7">{message.time}</div>
					</div>
					<div>{message.data}</div>
				</div>
			</div>
		</div>
	);
};

export default Message;
