import React from "react";
import { Link } from "react-router-dom";

export default class AignIn extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            email: "",
            password: ""
        };
    };

    render() {
        let imgStyle = {
            maxWidth : "200px"
        }

        return (
            <div className="container">
                <div className="col s12 m2" style={{ margin: "10px" }}>
                    <div className="z-depth-2 " style={{ padding: "10px" }}>

                        <h4>{this.props.summary.siteName}</h4>
                        <a href={this.props.summary.url}> visit site</a>
                        <h5>{this.props.summary.title}</h5>
                        <h6>{this.props.summary.description}</h6>
                        {
                            this.props.summary.images.map(image =>
                                <img style = {imgStyle} src={image.url} />
                            )
                        }
                    </div>
                </div>
            </div>
        );
    }
}
