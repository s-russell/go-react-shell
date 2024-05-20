import Routes from "./routes";
import {Link} from "wouter";
import {FC, useContext} from "react";
import {AuthContext} from "./api/auth-provider.tsx";
import {User} from "./api/user-svc.ts";


const Login: FC = () => <p>(<Link href="login">log in</Link>)</p>
const UserId: FC<{ user: User }> = ({user}) => <><p>(User: {user.name})</p></>

export default function RootLayout() {
    const {user} = useContext(AuthContext)

    return <>
        <h1>My Voltron Clusters</h1>
        {user == null ? <Login/> : <UserId user={user}/>}
        <hr/>
        <br/>
        <Routes/>
    </>
}