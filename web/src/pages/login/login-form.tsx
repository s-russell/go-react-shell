import { FC, FormEventHandler, useState } from "react"

export interface LoginFormData {
    username: string
    password: string
}



const LoginForm: FC<{ handleSubmit: (data: LoginFormData) => void }> = ({ handleSubmit }) => {
    const [data, setData] = useState({ username: "", password: "" })

    const handleFormSubmission: FormEventHandler<HTMLFormElement> = (event) => {
        event.preventDefault()
        handleSubmit(data)
    }
    return <>
        <form onSubmit={handleFormSubmission}>
            <fieldset>
                <legend>Enter you credentials:</legend>
                <label htmlFor="username">username</label>
                <input id="username" type="text"
                    value={data.username}
                    onChange={({ target: { value } }) => setData({ ...data, username: value })} />
                <label htmlFor="password">password</label>
                <input id="password" type="password"
                    value={data.password}
                    onChange={({ target: { value } }) => setData({ ...data, password: value })} />
                <button type="submit">Submit</button>
            </fieldset>
        </form>
    </>
}

export default LoginForm