import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";

export default function CreateAccountDialog() {
	return (
		<Dialog>
			<DialogTrigger asChild>
				<Button variant="outline">Войти</Button>
			</DialogTrigger>
			<DialogContent className="sm:max-w-[425px]">
				<Card>
					<CardHeader className="space-y-1">
						<CardTitle className="text-2xl">Войти</CardTitle>
						<CardDescription>
							Введите email и пароль
						</CardDescription>
					</CardHeader>
					<CardContent className="grid gap-4">
						<div className="grid gap-2">
							<Label htmlFor="email">Email</Label>
							<Input id="email" type="email" placeholder="m@example.com" />
						</div>
						<div className="grid gap-2">
							<Label htmlFor="password">Пароль</Label>
							<Input id="password" type="password" />
						</div>
					</CardContent>
					<CardFooter>
						<Button className="w-full">Войти</Button>
					</CardFooter>
				</Card>
			</DialogContent>
		</Dialog>
	);
}
