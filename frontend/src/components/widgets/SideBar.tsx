import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "../ui/button";
import { useParams } from "react-router-dom";

export default function SideBar() {
	const { id } = useParams<{ id: string }>();
	return (
		<ScrollArea className="h-[calc(100vh-88px)] w-[350px] border-r p-4 border-border-color">
			<Button className="w-[90%] p-2 bg-border-color mb-2">Новый чат</Button>
			{/* <Separator className="bg-border-color"/> */}
			<div className="my-4">Мои чаты</div>
			<Button variant="secondary" className="w-full p-2 mb-2" key={id}>{id}</Button>
		</ScrollArea>
	);
}
