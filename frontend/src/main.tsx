import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { BrowserRouter } from "react-router-dom";
import { AuthProvider } from "./context/Authprovider.tsx";

ReactDOM.createRoot(document.getElementById("root")!).render(
	<BrowserRouter>
		<React.StrictMode>
			<AuthProvider>
				<App />
			</AuthProvider>
		</React.StrictMode>
	</BrowserRouter>
);
