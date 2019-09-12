import React from "react";

export default class UserSearchResult extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
        };
        this.addUser = this.addUser.bind(this);

    };

    addUser(id) {
        console.log("clicked on user id: "+ id)
    }

    render() {
        return (
            <div >
            <div className = "row" style ={{margin : "15px 5px 15px 5px"}}>
            <img src={this.props.photo+ "?s=30"} alt="" className="circle responsive-img col" />
            <div className = "col"> {this.props.userName}</div>
            <div className = "col"> {this.props.firstName}</div>
            <div className = "col"> {this.props.lastName}</div>
            <div className = "col"> {this.props.name}</div>            
            </div>
          </div>
        );
    }
}
