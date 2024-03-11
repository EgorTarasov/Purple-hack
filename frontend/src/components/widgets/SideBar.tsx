import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "../ui/button";
import { useParams, useNavigate } from "react-router-dom";
import uuid from "react-uuid";
import { useToast } from "../ui/use-toast";
import { ToastAction } from "../ui/toast";

export default function SideBar() {
	const { id } = useParams<{ id: string }>();
	const navigate = useNavigate();

	const { toast } = useToast();

	function handleNewChat() {
		toast({
			title: "История чата не сохранится",
			description: "Для сохранения истории необходимо войти",
			action: <ToastAction altText="Новый чат" onClick={() => {navigate(`/chat/${uuid()}`);}}>Все равно создать</ToastAction>,
		});
	}

	return (
		<div className="border-r p-4 border-border-color w-[300px]">
			<Button className="w-[100%] mb-2 bg-border-color" onClick={handleNewChat}>
				Новый чат
			</Button>
			<ScrollArea className="h-[calc(100vh-200px)]">
				<div className="my-4">Мои чаты</div>
				<Button variant="secondary" className="max-w-full p-2 mb-2" key={id}>
					{id}
				</Button>
			</ScrollArea>
		</div>
	);
}
