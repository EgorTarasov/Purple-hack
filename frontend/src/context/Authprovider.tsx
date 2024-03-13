import { sessionToMessage } from "@/lib/sessionToMessage";
import storage from "@/lib/storage";
import { IMessage, ISession } from "@/models";
import ApiSession from "@/services/apiSession";
import {
	createContext,
	useContext,
	useState,
	ReactNode,
	useEffect,
} from "react";

interface AuthContextProps {
	isAuthorized: boolean;
	setIsAuthorized: React.Dispatch<React.SetStateAction<boolean>>;
	userSessions: ISession[];
	setUserSessions: React.Dispatch<React.SetStateAction<ISession[]>>;
	messageHistoryLists: IMessage[][];
	setMessageHistoryLists: React.Dispatch<React.SetStateAction<IMessage[][]>>;
	messageHistoryListCurrent: IMessage[];
	setMessageHistoryListCurrent: React.Dispatch<
		React.SetStateAction<IMessage[]>
	>;
}

const AuthContext = createContext<AuthContextProps>({
	isAuthorized: false,
	setIsAuthorized: () => {},
	userSessions: [],
	setUserSessions: () => {},
	messageHistoryLists: [],
	setMessageHistoryLists: () => {},
	messageHistoryListCurrent: [],
	setMessageHistoryListCurrent: () => {},
});

// function getCookie(name: string): string {
// 	const nameLenPlus = name.length + 1;
// 	return document.cookie
// 		.split(";")
// 		.map((c) => c.trim())
// 		.filter((cookie) => {
// 			return cookie.substring(0, nameLenPlus) === `${name}=`;
// 		})[0];
// }


export function AuthProvider({ children }: { children: ReactNode }) {
	const [isAuthorized, setIsAuthorized] = useState<boolean>(storage.getToken());
	const [userSessions, setUserSessions] = useState<ISession[]>([]);
	const [messageHistoryLists, setMessageHistoryLists] = useState<IMessage[][]>([]);
	const [messageHistoryListCurrent, setMessageHistoryListCurrent] = useState<IMessage[]>([]);

	const value = {
		isAuthorized,
		setIsAuthorized,
		userSessions,
		setUserSessions,
		messageHistoryLists,
		setMessageHistoryLists,
		messageHistoryListCurrent,
		setMessageHistoryListCurrent,
	};

	// useEffect(() => {
	// 	setIsAuthorized(getCookie("auth") !== undefined);
	// }, []);

	useEffect(() => {
		async function fetchData() {
			if (isAuthorized) {
				try {
					// Получение сессий
					const sessions = await ApiSession.getUserSession();
					// setUserSessions(sessions.data);
					setUserSessions(sessions.data);

					console.log("sessions", userSessions);
					console.log("sessions from resp", sessions);

					if (sessions.data) {
						const newMessageLists: IMessage[][] = sessionToMessage(sessions.data);
						setMessageHistoryLists(newMessageLists);
					}
				} catch (error) {
					console.log(error);
				}
			}
		}

		fetchData();
		console.log('here2')

		console.log("cookie", document.cookie);
		console.log("sessions", userSessions);
	}, [isAuthorized]);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
	return useContext(AuthContext);
}
