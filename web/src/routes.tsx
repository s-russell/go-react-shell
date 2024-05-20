import {Route, Switch} from "wouter";
import App from "./App";
import Login from "./pages/login/login";

const Routes = () => <>
    <Switch>
        <Route path="/">
            <App/>
        </Route>
        <Route path="/login">
            <Login/>
        </Route>
    </Switch>
</>

export default Routes