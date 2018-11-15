import React from "react";
import {Route, Switch} from "react-router-dom";
import Home from "./containers/Home";
import NotFound from "./containers/NotFound";
import Admin from "./containers/Admin";
import Order from "./containers/Order"

export default () => 
    <Switch>
        <Route path="/" exact component={Home} />
        <Route path="/admin" exact component={Admin} />
        <Route path="/orders*" exact component={Order} />
        <Route component={NotFound} />
    </Switch>