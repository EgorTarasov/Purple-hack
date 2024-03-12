import * as z from "zod"


export const SigninValidationSchema = z.object({
    email: z.string().email(),
    password: z.string(),
})