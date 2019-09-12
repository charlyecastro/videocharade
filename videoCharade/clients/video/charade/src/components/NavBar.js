
import React from "react";

import constants from "./Constants";
import { Navbar, NavItem } from 'react-materialize'


export default class NavBar extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            token: window.localStorage.getItem("token"),
        };
    };

    handleSubmit(evt) {
        evt.preventDefault();
    }

    deleteSession() {
        console.log("logging out!")
        fetch("https://"+constants.API+"/v1/sessions/mine",
            {
                method: "DELETE",
                headers: new Headers({
                    'Authorization': this.state.token,
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                localStorage.removeItem("token");
                this.setState({ user: {} })
                this.props.history.push(constants.routes.signin)
                return res;
            }).catch(err => {
                console.log(err)
                localStorage.removeItem("token");
                // return err.text().then(msg => {
                //     localStorage.removeItem("token");
                //     this.setState({
                //         errorMessage: msg.toString(),
                //         status: err.status
                //     })
                // })
            })
    }

    render() {
        const imgStyle = {
            paddingLeft: "10px"
        };

        return (

                <Navbar className = "teal lighten-1" >
                    <ul className="left">
                        <li> <img style={imgStyle} src={this.props.photoURL + "?s=60"} alt="" className="circle responsive-img" />   </li>
                        <li> <h5 style={{ paddingLeft: "10px", paddingTop: "10px" }}>{this.props.userName}</h5></li>
                    </ul>
                    <NavItem className="right"onClick={() => this.deleteSession()}>Sign out</NavItem>
                    <NavItem className="right"onClick={() => this.props.history.push((constants.routes.update))}>Update </NavItem>
                    <NavItem className="right"onClick={() => this.props.history.push((constants.routes.messages))}>Conversations</NavItem>
                    <NavItem className="right"onClick={() => this.props.history.push((constants.routes.leaderboard))}>Leaderboard</NavItem>
                    <NavItem className="right"onClick={() => this.props.history.push((constants.routes.charade))}>Charades</NavItem>
                </Navbar>

        );
    }
}