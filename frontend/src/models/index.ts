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
    error: boolean;
    sessionUuid?: string;
}

export interface ISession {
    id: string;
    queries: IQuery[];
    responses: IResponses[];
    createdAt: string;
}

export interface IQuery {
    id: number;
    model: string;
    body: string;
    createdAt: string;
}

export interface IResponses {
    id: number;
    context: { [key: string]: string[] };
    body: string;
    createdAt: string;
}