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

export interface ISession {
    id: string;
    queries: IQuery[];
    responses: IResponses[];
    createdAt: Date;
}

export interface IQuery {
    id: number;
    model: string;
    body: string;
    createdAt: Date;
}

export interface IResponses {
    id: number;
    context: {string: string[]};
    body: string;
    createdAt: Date;
}