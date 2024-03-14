const storagePrefix = "cbr";

const storage = {
    getToken: () => {
        return JSON.parse(
            window.localStorage.getItem(`${storagePrefix}token`) as string,
        );
    },
    setToken: (token: string) => {
        window.localStorage.setItem(
            `${storagePrefix}token`,
            JSON.stringify(token),
        );
    },
    clearToken: () => {
        window.localStorage.removeItem(`${storagePrefix}socket`);
    },
    getSocketUuid: () => {
        return JSON.parse(
            window.localStorage.getItem(`${storagePrefix}socket`) as string,
        );
    },
    setSocketUuid: (socket: string) => {
        window.localStorage.setItem(
            `${storagePrefix}socket`,
            JSON.stringify(socket),
        );
    },
    clearSocketUuid: () => {
        window.localStorage.removeItem(`${storagePrefix}socket`);
    },
};

export default storage;
