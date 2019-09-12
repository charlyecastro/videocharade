import React from "react";

import constants from "./Constants";

export default class Message extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            chanID : "",
            messID : "",
        };
    };

    editMessage(msgID) {

      let newMessage =  prompt("Please enter your new message")
      console.log(newMessage)
      let data = { "body": newMessage }
      console.log(data)
      data = JSON.stringify(data)

      if (newMessage != null) {

          fetch("https://"+ constants.API +"/v1/messages/" + msgID,
            {
                method: "PATCH",
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
                return res.json();
            }).then(json => {
                console.log(json)
            }).catch(err => {
                console.log(err)
            })
      }
    }

    deleteMessage(msgID) {
          fetch("https://" + constants.API + "/v1/messages/" + msgID,
            {
                method: "DELETE",
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': window.localStorage.getItem("token"),
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                console.log(res)
                //return res.json();
            }).then(json => {
                console.log(json)
            }).catch(err => {
                console.log(err)
            })
    }

    render() {

        let msgStyle = {
            borderRadius: "5px",
            backgroundColor: "blue",
            color: "white",
            padding: "10px",
            display: "block",
            margin: "20px",
            // position : "relative"
        }

        return (
            <div className="z-depth-1 light-blue " style={msgStyle}>
           <div className = "row">
           <div className = "col" style = {{fontWeight : "bold", marginBottom: "0px"}}>{this.props.name}</div>
           {this.props.userID == this.props.creator ?
                             <div style = {{marginTop : "0px", paddingTop : "0px"}}>
                             <a style = {{color : "white"}} onClick={() => this.editMessage(this.props.messID)} className = "col right">edit </a>
                             <a style = {{color : "white"}}  onClick={() => this.deleteMessage(this.props.messID)} className = "col right">delete</a>
                             </div>
                              : undefined
                        }
          
           </div>
            
            {this.props.body}
            
            </div>
        );
    }
}
