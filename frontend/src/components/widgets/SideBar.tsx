import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "../ui/button";

export default function SideBar() {
	return (
		<ScrollArea className="h-[calc(100vh-88px)] w-[350px] border-r p-4 border-border-color">
			<Button className="w-full p-2 bg-border-color mb-2">Новый чат</Button>
            <div>
                Мои чаты
            </div>
		</ScrollArea>
	);
}
