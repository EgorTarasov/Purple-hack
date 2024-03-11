import { IMessage } from "@/models";
import {
	createContext,
	Dispatch,
	ReactNode,
	SetStateAction,
	useContext,
	useEffect,
	useRef,
	useState,
} from "react";

export type WebsocketContextType = [
	IMessage[],
	Dispatch<SetStateAction<IMessage[]>>,
	boolean,
	string,
	(data: string) => void
];

export const WebsocketContext = createContext<WebsocketContextType>([
	[],
	() => {},
	false,
	"",
	() => {},
]);

interface WebsocketProviderProps {
	children: ReactNode;
	socketUuid: string;
	messageListDefault: IMessage[];
	modelType: string;
}

export const WebsocketProvider = ({
	children,
	socketUuid,
	messageListDefault,
	modelType,
}: WebsocketProviderProps) => {
	const [isReady, setIsReady] = useState<boolean>(false);
	const [val, setVal] = useState<string>("");
	const ws = useRef<WebSocket | null>(null);
	const [messageList, setMessageList] =
		useState<IMessage[]>(messageListDefault);

	useEffect(() => {
		setMessageList([]);
		const socket = new WebSocket(
			`wss://echo.websocket.events/${socketUuid}?model=${modelType}`
		);

		socket.onopen = () => setIsReady(true);
		socket.onclose = () => setIsReady(false);
		socket.onmessage = (event) => setVal(event.data);

		ws.current = socket;

		return () => {
			if (ws.current) {
				ws.current.close();
			}
		};
	}, [socketUuid, modelType]);

	const sendMessage = (data: string) => {
		if (ws.current) {
			ws.current.send(data);
		}
	};

	return (
		<WebsocketContext.Provider
			value={[messageList, setMessageList, isReady, val, sendMessage]}
		>
			{children}
		</WebsocketContext.Provider>
	);
};

// eslint-disable-next-line react-refresh/only-export-components
export function useWS() {
	return useContext(WebsocketContext);
}