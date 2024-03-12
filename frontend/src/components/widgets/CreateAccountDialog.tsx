import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "../ui/card";
import { useAuth } from "@/context/Authprovider";
import ApiAuth from "@/services/apiAuth";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "../ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { SigninValidationSchema } from "@/lib/validation";
import { z } from "zod";
import { toast } from "../ui/use-toast";
import PropagateLoader from "react-spinners/PropagateLoader";
import { useState } from "react";

export default function CreateAccountDialog() {
	const { isAuthorized, setIsAuthorized } = useAuth();
	const [isUserLoading, setIsUserLoading] = useState(false);

	// 1. Define your form.
	const form = useForm<z.infer<typeof SigninValidationSchema>>({
		resolver: zodResolver(SigninValidationSchema),
		defaultValues: {
			email: "",
			password: "",
		},
	});

	// 2. Define a submit handler.
	async function onSubmit(values: z.infer<typeof SigninValidationSchema>) {
		setIsUserLoading(true);
		try {
			const session = await ApiAuth.loginUser({
				email: values.email,
				password: values.password,
			});

			console.log("resp", session);
			form.reset();
			setIsAuthorized(true);
			return true;
		} catch (error) {
			return toast({
				title: "Ошибка авторизации. Попробуйте снова",
				variant: "destructive",
			});
		} finally {
			setIsUserLoading(false);
		}
	}

	return (
		<>
			{!isAuthorized && (
				<Dialog>
					<DialogTrigger asChild>
						<Button variant="outline">Войти</Button>
					</DialogTrigger>
					<DialogContent className="sm:max-w-[425px]">
						<Card className="shadow-lg">
							<CardHeader className="space-y-1">
								<CardTitle className="text-2xl">Войти</CardTitle>
								<CardDescription>Введите email и пароль</CardDescription>
							</CardHeader>
							<CardContent className="grid gap-4">
								<Form {...form}>
									<div className="sm:w-420 flex-center flex-col">
										<form
											onSubmit={form.handleSubmit(onSubmit)}
											className="flex flex-col gap-5 w-full mt-4 md:max-w-96"
										>
											<FormField
												control={form.control}
												name="email"
												render={({ field }) => (
													<FormItem>
														<FormLabel>Email</FormLabel>
														<FormControl>
															<Input
																type="email"
																className="shad-input"
																{...field}
															/>
														</FormControl>
														<FormMessage className="shad-form_message" />
													</FormItem>
												)}
											/>
											<FormField
												control={form.control}
												name="password"
												render={({ field }) => (
													<FormItem>
														<FormLabel>Пароль</FormLabel>
														<FormControl>
															<Input
																type="password"
																className="shad-input"
																{...field}
															/>
														</FormControl>
														<FormMessage className="shad-form_message" />
													</FormItem>
												)}
											/>
											<Button type="submit" className="shad-button_primary">
												{isUserLoading ? (
													<PropagateLoader color="#E3E8FD" size={5} />
												) : (
													"Войти"
												)}
											</Button>
										</form>
									</div>
								</Form>
							</CardContent>
						</Card>
					</DialogContent>
				</Dialog>
			)}
		</>
	);
}
