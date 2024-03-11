import { createContext, useContext, useState, ReactNode } from "react";

interface AuthContextProps {
	isAuthorized: boolean;
	setIsAuthorized: React.Dispatch<React.SetStateAction<boolean>>;
}

const AuthContext = createContext<AuthContextProps>({
	isAuthorized: getCookie("auth") !== "",
	setIsAuthorized: () => {},
});

function getCookie(name: string): string {
	const nameLenPlus = name.length + 1;
	return (
		document.cookie
			.split(";")
			.map((c) => c.trim())
			.filter((cookie) => {
				return cookie.substring(0, nameLenPlus) === `${name}=`;
			})
			.map((cookie) => {
				return decodeURIComponent(cookie.substring(nameLenPlus));
			})[0]
	);
}
export function AuthProvider({ children }: { children: ReactNode }) {
	const [isAuthorized, setIsAuthorized] = useState<boolean>(getCookie("auth") !== "");
	const value = { isAuthorized, setIsAuthorized };

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
	return useContext(AuthContext);
}
