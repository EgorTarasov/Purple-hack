export interface INavLink {
    route: string;
    label: string;
}

export interface IBubblePath {
    path: string;
    label: string;
    marginTop: number;
}

export interface IMessage {
    id: string;
    senderChat: boolean; 
    data: string;
    time: string;
}