import {Route, Switch} from "wouter";
import App from "./App";
import Login from "./pages/login/login";
import Profile from "./pages/profile/profile.tsx";

const Routes = () => <>
    <Switch>
        <Route path="/">
            <App/>
        </Route>
        <Route path="/login">
            <Login/>
        </Route>
        <Route path="/profile">
            <Profile/>
        </Route>
    </Switch>
</>

export default Routes