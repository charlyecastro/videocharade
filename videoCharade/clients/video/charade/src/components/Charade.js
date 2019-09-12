import React from "react";

import constants from "./Constants";
import NavBar from "./NavBar.js";
import UserSearch from "./UserSearch";
import Invitation from "./Invitation";
import Clock from "./Clock.js";
import { Navbar, NavItem } from 'react-materialize'

let myPeerConnection
export default class Charade extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            addedUser: false,
            inviteDeclined: false,
            isGameInSession: false,
            isVideoInSession: false,
            isGameOver: false,
            correctGuess: false,
            wasSkipped: false,
            newGuess: "",
            currWord: "",
            guessWord: "",
            numRight: 0,
            numSkipped: 0,
            numGuessed: 0,
            user: {},
            target: null, // this is for testing purposes i will set it to selected user
            isCaller: false,
            invitation: false,
            inviterName: "",
            videoOffer: null,
            accepted: false,
            token: window.localStorage.getItem("token"),
            hasMedia: false,
            localStream: null,
            myPeerConnection: null,
            peerConfig: { 'iceServers': [{ 'url': 'stun:stun.services.mozilla.com' }, { 'url': 'stun:stun.l.google.com:19302' }] },
            mediaConstraints: {
                audio: false, // We want an audio track
                video: true // ...and we want a video track
            },
        };

    };

    componentWillMount() {
        let sock = new WebSocket(`wss://${constants.API}/ws?auth=${window.localStorage.getItem("token")}`)

        sock.onopen = () => {
            console.log("Socket opened");
        };
        sock.onclose = () => {
            console.log("Socket closed");
        }

        sock.onmessage = (message) => {
            console.log("onmessage called!!")
            let msg = JSON.parse(message.data)
            switch (msg.type) {
                // Signaling messages: these messages are used to trade WebRTC
                // signaling information during negotiations leading up to a video
                // call.

                case "video-invitation":  // Invitation and offer to chat
                    this.setState({
                        invitation: true,
                        inviterName: msg.data.callerInfo.userName,
                        target: msg.data.callerInfo.id
                        //videoOffer: msg.data
                    })
                    console.log("hey we got a video invitation!!")
                    break;

                case "video-invitation-response":  // Invitation and offer to chat
                    if (msg.data.accepted) {
                        this.setState({ isVideoInSession: true })
                        this.beginCall()
                    } else {
                        this.setState({ inviteDeclined: true })
                    }
                    console.log("hey we got a video invitation response!!")
                    break;

                case "video-offer":  // Invitation and offer to chat
                    this.handleVideoOfferMsg(msg.data);
                    console.log("hey we got a video offer!!")
                    break;

                case "video-answer":  // Callee has answered our offer
                    this.handleVideoAnswerMsg(msg.data);
                    console.log("hey we got a video answer!!")
                    break;

                case "new-ice-candidate": // A new ICE candidate has been received
                    this.handleNewICECandidateMsg(msg.data);
                    console.log("hey we got a new ice candidate!!")
                    break;

                case "hang-up": // The other peer has hung up the call
                    this.closeVideoCall()
                    break;

                case "game-start":
                    //  //set up this!!
                    this.setState({ isGameInSession: true })
                    console.log("GAME STARTED")
                    this.renderGame(msg.data)
                    break;

                case "game-end":
                    //this.endGame() //set up this!!
                    console.log("GAME ENDED")
                    this.setState({ isGameOver: true })
                    break;

                case "guess":
                    let gameData = msg.data.GameState
                    console.log("THERE WAS A GUESS")
                    console.log(msg.data.correct)
                    this.setState({
                        guessWord: msg.data.guessed,
                        correctGuess: msg.data.correct,
                        wasSkipped: false,
                    })
                    this.renderGame(gameData)

                    break;

                case "skip":
                    this.setState({
                        wasSkipped: true,
                        correctGuess: false,
                    })
                    this.renderGame(msg.data)
                    break;

                default:
                    console.log("Unknown message received:");
            }

        }
        this.fetchAccount()
    }


    renderGame(data) {
        console.log("render GAME!")
        console.log(data)
        this.setState({
            currWord: data.currWord,
            numGuessed: data.numGuessed,
            numRight: data.numRight,
            numSkipped: data.numSkipped
        })
    }


    createPeerConnection() {
        myPeerConnection = new RTCPeerConnection(this.state.peerConfig)

        //sends iceCandidates to server when an iceCandidate is returned from the STUN server
        myPeerConnection.onicecandidate = (event) => {
            if (event.candidate) {
                console.log(event)
                let data = {
                    type: "new-ice-candidate",
                    data: { candidate: event.candidate },
                    userList: [this.state.target]
                }
                this.sendToServer(data)

            }
        }

        // myPeerConnection.ontrack = (event) => {
        //     this.remoteVideo.srcObject = event.streams[0]
        // }

        // myPeerConnection.onaddstream = (event) => {
        //     this.remoteVideo.srcObject = event.stream;
        // }

        if (myPeerConnection.addTrack !== undefined) {
            myPeerConnection.ontrack = (event) => {
                this.remoteVideo.srcObject = event.streams[0]
            }
        } else {
            myPeerConnection.onaddstream = (event) => {
                this.remoteVideo.srcObject = event.stream;
            }
        }


        myPeerConnection.onnegotiationneeded = () => {
            if (this.state.isCaller) {
                myPeerConnection.createOffer().then((offer) => {
                    return myPeerConnection.setLocalDescription(offer);
                })
                    .then(() => {
                        this.sendToServer({
                            type: "video-offer",
                            data: {
                                sdp: myPeerConnection.localDescription
                            },
                            userList: [this.state.target]
                        });
                    })
                    .catch(e => console.log("there was an error: ", e));
            }
        }
    }

    inviteUser() {
        let data = {
            type: "video-invitation",
            data: { callerInfo: this.state.user },
            userList: [this.state.target]
        }
        console.log("invitation!")
        console.log(data)
        this.sendToServer(data)
    }

    addUser = (userId) => {
        if (userId != this.state.user.id) {
            this.setState({
                addedUser: true,
                target: userId
            })
        }
    }

    respondToInvitation(response) {
        this.setState({ invitation: false })
        let data = {
            type: "video-invitation-response",
            data: {
                accepted: response
            },
            userList: [this.state.target]
        }
        this.sendToServer(data)
    }

    //Creates an offer for the other user r
    beginCall() {

        if (myPeerConnection) {
            alert("You can't start a call because you already have one open!");
        } else {
            this.setState({ isCaller: true })
            this.createPeerConnection();

            navigator.mediaDevices.getUserMedia(this.state.mediaConstraints)
                .then((stream) => {
                    // document.getElementById("local_video").src = window.URL.createObjectURL(localStream);
                    // document.getElementById("local_video").srcObject = localStream;
                   // this.myVideo.src = window.URL.createObjectURL(stream)
                    this.myVideo.srcObject = stream
                    if (myPeerConnection.addTrack !== undefined) {
                        stream.getTracks().forEach(track => myPeerConnection.addTrack(track, stream));
                    } else {
                        myPeerConnection.addStream(stream);
                    }


                    // stream.getTracks().forEach(track => myPeerConnection.addTrack(track, stream));
                })
                .catch(e => console.log("there was an issue getting media: " + e));
        }
    }

    hangUp() {
        this.closeVideoCall();
        this.sendToServer({
            type: "hang-up",
            data: null,
            userList: [this.state.target],
        });
    }

    closeVideoCall() {
        // Close the RTCPeerConnection
        if (myPeerConnection) {

            this.setState({ isCaller: false, isVideoInSession: false, isGameInSession: false })
            // Disconnect all our event listeners; we don't want stray events
            // to interfere with the hangup while it's ongoing.
            myPeerConnection.onaddstream = null;  // For older implementations
            myPeerConnection.ontrack = null;      // For newer ones
            myPeerConnection.onnicecandidate = null;
            myPeerConnection.onnotificationneeded = null;
            // Stop the videos
            if (this.myVideo.srcObject) {
                this.myVideo.srcObject.getTracks().forEach(track => track.stop());
            }
            if (this.remoteVideo.srcObject) {
                this.remoteVideo.srcObject.getTracks().forEach(track => track.stop());
            }
            this.myVideo.src = null;
            this.remoteVideo.src = null;
            // Close the peer connection
            myPeerConnection.close();
            myPeerConnection = null;
        }
        // Disable the hangup button
        //targetUsername = null;
    }

    // receive video offer from user that want to connecto to you, by sending a video-answer to the server
    handleVideoOfferMsg(msg) {
        this.setState({
            accepted: true,
            isVideoInSession: true
        })
        let localStream = null;

        this.createPeerConnection();
        let desc = new RTCSessionDescription(msg.sdp);

        myPeerConnection.setRemoteDescription(desc)
            .then(() => {
                return navigator.mediaDevices.getUserMedia(this.state.mediaConstraints);
            })
            .then((stream) => {
                localStream = stream;
                //this.myVideo.src = window.URL.createObjectURL(stream)
                this.myVideo.srcObject = stream;

                if (myPeerConnection.addTrack !== undefined) {
                    localStream.getTracks().forEach(track =>
                        myPeerConnection.addTrack(track, localStream)
                    );
                } else {
                    myPeerConnection.addStream(localStream);
                }


                //localStream.getTracks().forEach(track => myPeerConnection.addTrack(track, localStream));
            })
            .then(() => {
                return myPeerConnection.createAnswer();
            })
            .then((answer) => {
                return myPeerConnection.setLocalDescription(answer);
            })
            .then(() => {
                var msg = {
                    type: "video-answer",
                    data: { sdp: myPeerConnection.localDescription },
                    userList: [this.state.target]
                };

                this.sendToServer(msg);
            })
            .catch((e) => { console.log(e) });
    }

    handleVideoAnswerMsg(msg) {
        console.log("Call recipient has accepted our call");

        // Configure the remote description, which is the SDP payload
        // in our "video-answer" message.

        let desc = new RTCSessionDescription(msg.sdp);
        myPeerConnection.setRemoteDescription(desc).catch(e => { console.log(e) });
    }

    handleNewICECandidateMsg(msg) {
        let candidate = new RTCIceCandidate(msg.candidate);

        myPeerConnection.addIceCandidate(candidate)
            .catch(e => { console.log(e) });
    }

    createOfferError(error) {
        console.log("method was called!!" + error);
    }

    //sends video chat information to to messaging microService
    sendToServer(data) {
        data = JSON.stringify(data)
        fetch("https://" + constants.API + "/v1/handleOffer",
            {
                method: "POST",
                body: data,
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': window.localStorage.getItem("token"),
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                console.log(res.status)
                // return res.json();
            }).then(json => {
                //console.log(json)
            }).catch(err => {
                console.log(err)
            })
    }

    //grab user data from server
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

    beginGame() {
        this.setState({ isGameOver: false })
        let data = {
            "FirstUserID": this.state.user.id,
            "SecondUserID": this.state.target
        }

        data = JSON.stringify(data)
        fetch("https://" + constants.API + "/v1/charades",
            {
                method: "POST",
                body: data,
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': window.localStorage.getItem("token"),
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                console.log(res.status)
                return res.json();
            }).then(json => {
                console.log(json)
            }).catch(err => {
                console.log(err)
            })
    }

    postGuess(evt) {
        evt.preventDefault()
        let data = {
            "FirstUserID": this.state.target,
            "SecondUserID": this.state.user.id,
            "Guess": this.state.newGuess
        }

        data = JSON.stringify(data)
        console.log(data)
        fetch("https://" + constants.API + "/v1/charades/guess",
            {
                method: "POST",
                body: data,
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': window.localStorage.getItem("token"),
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                console.log(res.status)
                return res.json();
            }).then(json => {
                console.log(json)
            }).catch(err => {
                console.log(err)
            })
    }

    postSkip() {
        let data = {
            "FirstUserID": this.state.user.id,
            "SecondUserID": this.state.target,
        }

        data = JSON.stringify(data)
        console.log(data)
        fetch("https://" + constants.API + "/v1/charades/skip",
            {
                method: "POST",
                body: data,
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': window.localStorage.getItem("token"),
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                console.log(res.status)
                return res.json();
            }).then(json => {
                console.log(json)
            }).catch(err => {
                console.log(err)
            })
    }


    render() {

        let btnMargins = {
            marginLeft: '20px',
            marginRight: '20px'
        };

        return (
            <div>
                <NavBar photoURL={this.state.user.photoURL} userName={this.state.user.userName} history={this.props.history}></NavBar>

                {this.state.inviteDeclined && !this.state.isVideoInSession ?
                    <div>
                        <div className="col s12 " style={{ margin: "10px" }}>
                            <div className="z-depth-2 " style={{ padding: "10px" }}>

                                <h5> Sorry. Your Invitation was declined :( </h5>
                                <a className="col" onClick={() => {
                                    this.setState({
                                        inviteDeclined: false,
                                    })
                                }}> close </a>
                            </div>
                        </div>
                    </div>
                    : undefined
                }

                {this.state.addedUser && !this.state.isVideoInSession ?
                    <div>
                        <div className="col s12 " style={{ margin: "10px" }}>
                            <div className="z-depth-2 " style={{ padding: "10px" }}>

                                <h5>  User Was Added! </h5>
                                <a className="col" onClick={() => {
                                    this.setState({
                                        addedUser: false,
                                    })
                                }}> close </a>
                            </div>
                        </div>
                    </div>
                    : undefined
                }

                {this.state.isVideoInSession && this.state.wasSkipped ?
                    <div>
                        <div className="col s12 yellow darken-1" style={{ margin: "10px" }}>
                            <div className="z-depth-2 " style={{ padding: "10px" }}>
                                <h5>  Word was Skipped </h5>
                                <a className="col" onClick={() => {
                                    this.setState({
                                        wasSkipped: false,
                                    })
                                }}> close </a>
                            </div>
                        </div>
                    </div>
                    : undefined
                }

                {this.state.isVideoInSession && this.state.correctGuess ?
                    <div>
                        <div className="col s12 teal lighten-1" style={{ margin: "10px" }}>
                            <div className="z-depth-2 " style={{ padding: "10px" }}>

                                <h5>  Guess Was Correct! Guess was {this.state.guessWord} </h5>
                                <a className="col" style={{ color: "black" }} onClick={() => {
                                    this.setState({
                                        correctGuess: false,
                                    })
                                }}> close </a>
                            </div>
                        </div>
                    </div>
                    : undefined
                }

                {this.state.isGameOver && this.state.isVideoInSession ?
                    <div>
                        <div className="col s12 " style={{ margin: "10px" }}>
                            <div className="z-depth-2 " style={{ padding: "10px" }}>

                                <h5> Times Up! Check out the leaderboard to see if you made it to the top 10!</h5>
                                <div className="row">
                                    <button style={btnMargins} className="col btn" onClick={() => this.props.history.push((constants.routes.leaderboard))}> See Leaderboard</button>
                                </div>
                            </div>
                        </div>
                    </div>


                    // <button className="btn" onClick={() => this.props.history.push((constants.routes.leaderboard))}>Leaderboard</button>
                    : undefined}
                {this.state.invitation ?
                    <div>
                        <div className="col s12 " style={{ margin: "10px" }}>
                            <div className="z-depth-2 " style={{ padding: "10px" }}>

                                <h5> {this.state.inviterName} wants to chat with you</h5>
                                <div className="row">
                                    <button style={btnMargins} className="col btn" onClick={() => this.respondToInvitation(true)}> accept</button>
                                    <button style={btnMargins} className="col btn red darken-1" onClick={() => this.respondToInvitation(false)}> decline</button>
                                </div>
                            </div>
                        </div>
                    </div>

                    //<Invitation className="col" name={this.state.inviterName}  />
                    : undefined
                }

                {!this.state.isVideoInSession ?
                    <div className="container">
                        <h3>Welcome to VideoCharades {this.state.user.firstName}</h3>
                        <h5> 1. Search a user</h5>
                        <h5> 2. Select the user</h5>
                        <h5> 3. Hit Call</h5>
                        <h5> 4. Then Press Start!</h5>
                        <UserSearch token={this.state.token} onAddUser={this.addUser} />
                    </div>
                    :
                    undefined
                }


                <div className="container row" style={{ paddingTop: "50px" }}>

                    {this.state.isGameInSession && this.state.isCaller ?
                        <h3 className="center-align">Current Word: {this.state.currWord}</h3>
                        :
                        undefined
                    }
                    <div className="row">
                    {this.state.isVideoInSession  ?
                        <div className="col" >
                        <video className="center-align" style={{ width: "200px", margin: "0px", padding: "0px" }} className="col" id="localVideo" autoPlay ref={(ref) => { this.myVideo = ref }} />
                        <video className="center-align" style={{ margin: "0px", padding: "0px" }} id="remoteVideo" autoPlay ref={(ref) => { this.remoteVideo = ref }} />
                    </div>
                        :
                        undefined
                    }
                        {/* <div className="col" >
                            <video className="center-align" style={{ width: "200px", margin: "0px", padding: "0px" }} className="col" id="localVideo" autoPlay ref={(ref) => { this.myVideo = ref }} />
                            <video className="center-align" style={{ margin: "0px", padding: "0px" }} id="remoteVideo" autoPlay ref={(ref) => { this.remoteVideo = ref }} />
                        </div> */}
                        {this.state.isGameInSession ?
                            <di className="col">
                                {/* <h3>Time: </h3> */}
                                <Clock startCount={60}></Clock>
                                <p>Right: {this.state.numRight}</p>
                                <p>Guesses: {this.state.numGuessed}</p>
                                <p>Skipped: {this.state.numSkipped}</p>
                                <p>Guess: {this.state.guessWord}</p>
                            </di>
                            :
                            undefined
                        }


                    </div>
                    {this.state.isGameInSession && !this.state.isCaller ?
                        <form onSubmit={evt => this.postGuess(evt)} >
                            <div className="row center spacing">
                                <div className="input-field col s12 l10">
                                    <input type="text"
                                        className="form-control"
                                        placeholder="Guess the word"
                                        onInput={evt => this.setState({ newGuess: evt.target.value })} />
                                </div>
                                <div className="form-group col s2" style={{ marginTop: "20px" }}>
                                    <button type="submit" className="btn" > Submit </button>
                                </div>
                            </div>
                        </form>
                        : undefined
                    }

                    <div className="center-align">
                        {this.state.isCaller && this.state.isVideoInSession && !this.state.isGameInSession ?
                            <button style={btnMargins} className="btn" onClick={() => this.beginGame()}>Start</button> :
                            undefined
                        }

                        {!this.state.isVideoInSession ?
                            <button style={btnMargins} className="btn" onClick={() => this.inviteUser()} >Call</button> :
                            undefined
                        }

                        <button style={btnMargins} className="btn red darken-1" onClick={() => this.hangUp()} ref={(ref) => { this.hangUpBtn = ref }} disabled={!this.state.isVideoInSession}>Hang Up</button>

                        {this.state.isCaller && this.state.isGameInSession ?
                            <button style={btnMargins} className="btn yellow darken-1" onClick={() => this.postSkip()}>Skip</button> :
                            undefined
                        }
                    </div>
                </div>
            </div>
        );
    }
}