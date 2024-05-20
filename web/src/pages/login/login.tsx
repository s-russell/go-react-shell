import {FC, useContext} from "react"
import LoginForm, {LoginFormData} from "./login-form"
import {AuthContext} from "../../api/auth-provider"
import {useLocation} from "wouter";

const Login: FC = () => {

    const auth = useContext(AuthContext)
    const [_, setLocation] = useLocation()

    const handleSubmit = ({ username, password }: LoginFormData) => {
        auth.login(username, password, (user) => setLocation("/"))
    }
    return <>
        {auth.isAuthenticating ? <h2>Logging in...</h2> : <LoginForm handleSubmit={handleSubmit} />}
    </>
}

export default Login