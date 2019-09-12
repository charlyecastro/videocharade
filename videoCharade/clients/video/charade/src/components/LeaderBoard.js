import React from "react";

import constants from "./Constants";
import NavBar from "./NavBar.js";

export default class LeaderBoard extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            token: window.localStorage.getItem("token"),
            scoreList: [],
            user : {}
        };
    };

    componentWillMount() {
        this.fetchLeaderBoard()
        this.fetchAccount()
    }

    fetchLeaderBoard() {
        fetch("https://" + constants.API + "/v1/leaderboards",
            {
                method: "GET",
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': this.state.token,
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                return res.json();
            }).then(json => {
                console.log(json)
                if (json) {
                this.setState({ scoreList: json })
                }
            }).catch(err => {
                console.log(err)
            })
    }

    fetchAccount() {
        fetch("https://" + constants.API + "/v1/users/me",
            {
                method: "GET",
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': this.state.token,
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                return res.json();
            }).then(json => {
                console.log(json)
                this.setState({ user: json })
            }).catch(err => {
                console.log(err)
                // return err.text().then(msg => {
                //     this.setState({
                //         errorMessage: msg.toString(),
                //         status: err.status
                //     })
                // })
            })
    }


    render() {
        return (
            <div>
                <NavBar photoURL={this.state.user.photoURL} userName={this.state.user.userName} history={this.props.history}></NavBar>
                <div className="container">
                    <h3 className="center">Leaderboard</h3>
                    <h5 className="center">These following are the top 10 best Pairs in VideoCharades</h5>
                    <table className="highlight centered z-depth-2">
                        <thead>
                            <tr>
                                <th>Actor</th>
                                <th>Guesser</th>
                                <th>Words Played</th>
                                <th>Correct Guesses</th>
                            </tr>
                        </thead>

                        <tbody>

                            {
                                this.state.scoreList.map(score =>
                                    <tr>
                                        <td>{score.actorID}</td>
                                        <td>{score.guesserID}</td>
                                        <td>{score.numPlayed}</td>
                                        <td>{score.numRight}</td>
                                    </tr>
                                )
                            }
                        </tbody>
                    </table>
                </div>
            </div>
        );
    }
}