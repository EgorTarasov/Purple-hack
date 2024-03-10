import { createContext, ReactNode, useEffect, useRef, useState } from "react";

type WebsocketContextType = [
  boolean,
  string,
  ((data: string) => void)
];

export const WebsocketContext = createContext<WebsocketContextType>([
  false,
  "",
  () => {}, 
]);

interface WebsocketProviderProps {
  children: ReactNode;
  socketUuid: string;
}

export const WebsocketProvider = ({ children, socketUuid }: WebsocketProviderProps) => {
  const [isReady, setIsReady] = useState<boolean>(false);
  const [val, setVal] = useState<string>("");
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    const socket = new WebSocket(`wss://echo.websocket.events/${socketUuid}`);

    socket.onopen = () => setIsReady(true);
    socket.onclose = () => setIsReady(false);
    socket.onmessage = (event) => setVal(event.data);

    ws.current = socket;

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [socketUuid]);

  const sendMessage = (data: string) => {
    if (ws.current) {
      ws.current.send(data);
    }
  };

  return (
    <WebsocketContext.Provider value={[isReady, val, sendMessage]}>
      {children}
    </WebsocketContext.Provider>
  );
};
