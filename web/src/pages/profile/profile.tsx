import {useContext} from "react";
import {AuthContext} from "../../api/auth-provider.tsx";

export default function() {
    const {user} = useContext(AuthContext)

    return user && <>
        <h1>Hello, {user?.firstName}</h1>
        <h2>Roles:</h2>
        <ul>
            {user.roles.map(r => <li key={r}>{r}</li> )}
        </ul>
    </>
}