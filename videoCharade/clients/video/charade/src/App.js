
import React, { Component } from 'react';
//import { HashRouter as Router, Switch, Redirect, Route } from 'react-router-dom';
import { BrowserRouter as Router, Switch, Redirect, Route } from 'react-router-dom';


import constants from "./components/Constants.js";

import SignIn from "./components/SignIn.js";
import SignUp from "./components/SignUp.js";
import Messages from "./components/Messages";
import Update from "./components/Update";
import Summary from "./components/Summary";
import Charade from "./components/Charade"
import LeaderBoard from "./components/LeaderBoard"

class App extends Component {
  render() {
    return (
      <Router>
        <Switch>
          <Route exact path={constants.routes.signin} component={SignIn} />
          <Route exact path={constants.routes.signup} component={SignUp} />
          <Route exact path={constants.routes.messages} component={Messages} />
          <Route exact path={constants.routes.update} component={Update} />
          <Route exact path={constants.routes.summary} component={Summary} />
          <Route exact path={constants.routes.charade} component={Charade} />
          <Route exact path={constants.routes.leaderboard} component={LeaderBoard} />
          <Redirect to={constants.routes.signin} />
        </Switch>
      </Router>
    );
  }
}

export default App;