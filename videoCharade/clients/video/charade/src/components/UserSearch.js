import React from "react";
import UserSearchResult from "./UserSearchResult";
import Alert from "./Alert.js";
import constants from "./Constants"

export default class UserSearch extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            prefix: "",
            userList: [],
            errorMessage: "",
            status: 0,
        };
    };

    compare(a, b) {
        if (a.userName < b.userName)
            return -1;
        if (a.userName > b.userName)
            return 1;
        return 0;
    }


    handleChange(e) {
        let prefix = e.target.value
         fetch("https://"+constants.API+"/v1/users?q=" + prefix,
            {
                method: "GET",
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': this.props.token,
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                return res.json();
            }).then(json => {
                json.sort(this.compare);
                this.setState({ userList: json })
            }).catch(err => {
                return err.text().then(msg => {
                    let errorList = [{id: 0, 
                        userName: "status " + err.status, 
                        firstName: msg.toString(), 
                        lastName: "",
                         photoURL: ""}]
                    // this.setState({
                    //     userList : errorList
                    // })
                })
            })
    }

    userClick(userId){
        console.log("you have been clicked")
        console.log(userId)
        this.setState({userList : []})
        this.props.onAddUser(userId); 
    }

    

    render() {
        return (
            <div>
                {this.state.errorMessage ?
                    <Alert msg={this.state.errorMessage} status={this.state.status} /> : undefined
                }
                <div className="row">
                    <div className="col s12">
                        <div className="row">
                            <div className="input-field col s12" style={{ position: "relative" }}>
                                <input type="text" id="userSearch" className="" 
                                onInput={(e) => { this.handleChange(e) }} style={{ marginBottom: '0px' }} />
                                <label htmlFor="autocomplete-input">Search Users</label>
                                <div className="z-depth-2 " style={{ position: "absolute", backgroundColor: "white", width: "98%" }}>
                                    {

                                        this.state.userList.map(user =>
                                            <div>

                                                {/* onClick={() => this.handleSelected(channel._id)} */}
                                                <button onClick={() => this.userClick(user.id)}>
                                            <UserSearchResult   photo={user.photoURL} userName={user.userName} firstName={user.firstName} lastName={user.lastName} key={user.id} />
                                            </button>
                                            </div>
                                        )
                                    }
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}
