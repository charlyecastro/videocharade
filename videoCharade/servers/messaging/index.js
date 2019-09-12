//put interpreter into strict mode
"use strict";

//require the server libraries
const express = require("express");
const morgan = require("morgan");
const bodyParser = require("body-parser")
//import mongoose
const mongoose = require("mongoose")
//import rabbit
const amqp = require('amqplib/callback_api');
const addr = process.env.MESSAGESADDR || ":80";
const [host, port] = addr.split(":");
const app = express();
const server = require('http').createServer(app);

//set up middlewares
const authHandler = require("./middleware/authenticate.js")
app.use(express.json());
app.use(morgan("dev"));
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());
app.use(authHandler)

//set up mongoDB
const Channel = require("./models/channel.js")
const Message = require("./models/message.js")
const mongoAddr = process.env.MONGOADDR || "localhost:27017";
const url = "mongodb://" + mongoAddr + "/messages"

//set up queue
const rabbitAddr = process.env.RABBITADDR;
const queueName = 'Messaging';
let messageChannel;

//function that creates a new Channel
function newChannel(body, user) {
    return new Channel({
        _id: new mongoose.Types.ObjectId(),
        name: body.name,
        nameLower: body.name.toLowerCase(),
        description: body.description,
        private: body.private,
        members: body.members,
        createdAt: new Date().getTime(),
        creator: user,
        editedAt: new Date().getTime(),
    })
}

//function that creates a new Message
function newMessage(body, chanID, user) {
    return new Message({
        _id: new mongoose.Types.ObjectId(),
        channelID: chanID,
        name: user.name,
        body: body.body,
        createdAt: new Date().getTime(),
        creator: user,
        editedAt: new Date().getTime(),
    })
}

// set up default General Channel once
const genBody = { name: "General", description: "This channel is for general convos", creator: { id: 0, name: "slack" } }

let name = genBody.name.toLowerCase()
Channel.findOne({ "nameLower": name }).exec()
    .then(doc => {
        if (!doc) {
            let general = newChannel(genBody, genBody.creator)
            general.save()
                .then(result => { console.log("addedGeneral: " + result) })
                .catch(err => { console.log(err) })
        } else {
            console.log("didnt add general")
        }
    })
    .catch(err => { console.log(err) })


// /v1/channels: refers to all channels
//-----------------------------------------------
// GET: Respond with the list of all channels (just the channel models, not the messages in those channels) that the current user is allowed to see, encoded as a JSON array. Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.
app.get("/v1/channels", (req, res) => {
    Channel.find().exec()
        .then(docs => {
            res.set("Content-Type", "application/json");
            res.json({ docs });
        })
        .catch(err => {
            console.log(err)
            res.sendStatus(500).end()
        })
});


// POST: Create a new channel using the channel model JSON in the request body. The name property is required, but description is optional. Respond with a 201 status code, a Content-Type set to application/json, and a copy of the new channel model (including its new ID) encoded as a JSON object.
app.post("/v1/channels", (req, res) => {
    let fetchedUser = req.get("X-User");
    let user = JSON.parse(fetchedUser);

    let name = req.body.name.toLowerCase()
    Channel.findOne({ "nameLower": name }).exec()
        .then(doc => {
            if (!doc) {
                let channel = newChannel(req.body, user)
                if (channel.private) {
                    channel.members.push(user)
                }
                channel.save()
                    .then(result => {
                        const newChannelMsg = {
                            "type": "channel-new",
                            "data": result,
                            "userList": []
                        };

                        console.log(newChannelMsg)
                        messageChannel.sendToQueue(
                            queueName,
                            new Buffer(JSON.stringify(newChannelMsg)),
                            { persistent: true }  // make sure messages are stored until ack'ed
                        );

                        res.set("Content-Type", "application/json");
                        res.status(201).json(result);

                    })
                    .catch(err => {
                        console.log(err)
                        res.sendStatus(500).end()
                    })
            } else {
                res.status(400).send("Channel Name is already in use")
            }
        })
        .catch(err => {
            console.log(err)
            res.sendStatus(500).end()
        })
});


// /v1/channels/{channelID}: refers to a specific channel identified by {channelID}
//--------------------------------------------------------------
// GET: If this is a private channel and the current user is not a member, respond with a 403 (Forbidden) status code. Otherwise, respond with the most recent 100 messages posted to the specified channel, encoded as a JSON array of message model objects. Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.
app.get("/v1/channels/:channelID", (req, res) => {
    let fetchedUser = req.get("X-User");
    let user = JSON.parse(fetchedUser);
    let queryID = req.query.id;
    let chanID = req.params.channelID
    Channel.findById(chanID).exec()
        .then(channel => {
            if (channel.private && channel.creator.id != user.id) {
                let memberList = channel.members
                let found = false
                for (let i = 0; i < memberList.length; i++) {
                    if (memberList[i].id == user.id) {
                        found = true;
                        break
                    }
                }
                if (!found) {
                    res.sendStatus(403)
                    res.end()
                    return
                }
            }
            let dbQuery;
            if (!queryID) {
                dbQuery = Message.find({ channelID: chanID })
            } else {
                dbQuery = Message.find({ _id: { $lt: queryID } })
            }
            dbQuery.limit(100).exec()
                .then(messages => {
                    res.set("Content-Type", "application/json");
                    res.json({ messages });
                })
                .catch(err => { 
                    console.log(err)
                    res.sendStatus(500).end() })

        }).catch(err => { 
            console.log(err)
            res.sendStatus(500).end() })
});

// POST: If this is a private channel and the current user is not a member, respond with a 403 (Forbidden) status code. Otherwise, create a new message in this channel using the JSON in the request body. The only message property you should read from the request is body. Set the others based on context. Respond with a 201 status code, a Content-Type set to application/json, and a copy of the new message model (including its new ID) encoded as a JSON object.
app.post("/v1/channels/:channelID", (req, res) => {
    let fetchedUser = req.get("X-User");
    let user = JSON.parse(fetchedUser);
    let chanID = req.params.channelID
    Channel.findById(chanID).exec()
        .then(doc => {
            if (doc.private && doc.creator.id != user.id) {
                let memberList = doc.members
                let found = false
                for (let i = 0; i < memberList.length; i++) {
                    if (memberList[i].id == user.id) {
                        console.log("member FOUND!!")
                        found = true;
                        break
                    }
                }
                if (!found) {
                    res.sendStatus(403).end()
                    return
                }
            }
            let message = newMessage(req.body, chanID, user)
            message.save()
                .then(result => {

                    let idList = doc.members.map(m => m.id);
                    let messData = JSON.stringify(result)
                    //console.log(messData)

                    const newMessageMsg = {
                        "type": "message-new",
                        "data": result,
                        "userList": idList,
                    };

                    console.log(newMessageMsg)
                    messageChannel.sendToQueue(
                        queueName,
                        Buffer.from(JSON.stringify(newMessageMsg)),
                        { persistent: true }  // make sure messages are stored until ack'ed
                    );
                    res.set("Content-Type", "application/json");
                    res.status(201).json(result);
                })
                .catch(err => { 
                    console.log(err)
                    res.sendStatus(500).end() })
        })
        .catch(err => {
            console.log(err)
            res.sendStatus(500).end()
        })
});

// PATCH: If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden). Otherwise, update only the name and/or description using the JSON in the request body and respond with a copy of the newly-updated channel, encoded as a JSON object. Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.

app.patch("/v1/channels/:channelID", (req, res) => {
    let fetchedUser = req.get("X-User");
    let user = JSON.parse(fetchedUser);
    let chanID = req.params.channelID
    Channel.findById(chanID).exec()
        .then(doc => {
            if (doc.creator.id != user.id) {
                res.sendStatus(403).end()
            } else {
                doc.name = req.body.name
                doc.nameLower = req.body.name.toLowerCase()
                doc.description = req.body.description || doc.description
                doc.editedAt = new Date().getTime(),
                    doc.save()
                        .then(result => {
                            const updateChannelMsg = {
                                "type": "channel-update",
                                "data": result,
                                "userList": []
                            };
                            console.log(updateChannelMsg)
                            messageChannel.sendToQueue(
                                queueName,
                                new Buffer(JSON.stringify(updateChannelMsg)),
                                { persistent: true }  // make sure messages are stored until ack'ed
                            );

                            res.set("Content-Type", "application/json");
                            res.json(result);
                        })
                        .catch(err => { 
                            console.log(err)
                            res.sendStatus(500).end() })
            }
        })
        .catch(err => {
            console.log(err)
            res.sendStatus(500).end()
        })
});


// DELETE: If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden). Otherwise, delete the channel and all messages related to it. Respond with a plain text message indicating that the delete was successful.
app.delete("/v1/channels/:channelID", (req, res) => {
    let fetchedUser = req.get("X-User");
    let user = JSON.parse(fetchedUser);
    let chanID = req.params.channelID
    Channel.findById(chanID).exec()
        .then(doc => {
            if (doc.creator.id != user.id) {
                res.sendStatus(403).end()
            } else {
                Channel.deleteOne({ _id: chanID }, err => {
                    if (err) {
                        res.sendStatus(500).end()
                    }
                });
                Message.deleteMany({ channelID: chanID }, err => {
                    if (err) {
                        res.sendStatus(500).end()
                    }
                });

                let idList = doc.members.map(m => m.id);
                const delChannelMsg = {
                    "type": "channel-delete",
                    "data": chanID,
                    "userList": []
                };
                console.log(delChannelMsg)
                messageChannel.sendToQueue(
                    queueName,
                    new Buffer(JSON.stringify(delChannelMsg)),
                    { persistent: true }  // make sure messages are stored until ack'ed
                );

                res.send("Delete was successful!")
            }
        })
        .catch(err => {
            console.log(err)
            res.sendStatus(500).end()
        })
});


// /v1/channels/{channelID}/members: refers to the members of a private channel identified by {channelID}
//--------------------------------------------------------------
// POST: If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden). Otherwise, add the user supplied in the request body as a member of this channel, and respond with a 201 status code and a simple plain text message indicating that the user was added as a member. Only the id property of the user is required, but the client may post the entire user profile.
app.post("/v1/channels/:channelID/members", (req, res) => {
    let fetchedUser = req.get("X-User");
    let user = JSON.parse(fetchedUser);
    let chanID = req.params.channelID
    Channel.findById(chanID).exec()
        .then(doc => {
            if (doc.creator.id != user.id) {
                res.sendStatus(403).end()
            } else {
                doc.members.push(req.body.member)
                doc.save()
                    .then(result => {
                        res.status(201).send("User was added");
                    })
                    .catch(err => { 
                        console.log(err)
                        res.sendStatus(500).end() })
            }
        })
        .catch(err => {
            console.log(err)
            res.sendStatus(500).end()
        })
});

app.delete("/v1/channels/:channelID/members", (req, res) => {
    let fetchedUser = req.get("X-User");
    let user = JSON.parse(fetchedUser);
    let chanID = req.params.channelID
    Channel.findById(chanID).exec()
        .then(doc => {
            if (doc.creator.id != user.id) {
                res.sendStatus(403).end()
            } else {
                Channel.findOneAndUpdate({ _id: chanID },
                    { $pull: { "members": { "id": req.body.id } } }).exec()
                    .then(result => {
                        res.status(200).send("User was removed");
                    })
                    .catch(err => {
                        console.log(err)
                        res.sendStatus(500).end()
                    })
            }
        })
        .catch(err => {
            console.log(err)
            res.sendStatus(500).end()
        })
});

// /v1/messages/{messageID}: refers to a specific message identified by {messageID}
//--------------------------------------------------------------
// PATCH: If the current user isn't the creator of this message, respond with the status code 403 (Forbidden). Otherwise, update the message body property using the JSON in the request body, and respond with a copy of the newly-updated message, encoded as a JSON object. Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.
app.patch("/v1/messages/:messageID", (req, res) => {
    let fetchedUser = req.get("X-User");
    let user = JSON.parse(fetchedUser);
    let messID = req.params.messageID
    Message.findById(messID).exec()
        .then(message => {
            if (message.creator.id != user.id) {
                res.sendStatus(403).end()
            } else {
                message.body = req.body.body
                message.editedAt = new Date().getTime()

                let chanID = message.channelID
                let members = []
                Channel.findById(chanID).exec()
                    .then(doc => {
                        members = doc.members
                    })
                    .catch(err => {
                        console.log(err)
                        res.sendStatus(500).end()
                        return
                    })

                message.save()
                    .then(newMessage => {
                        let idList = members.map(m => m.id);
                        const updateMessageMsg = {
                            "type": "message-update",
                            "data": newMessage,
                            "userList": idList,
                        };

                        console.log(updateMessageMsg)
                        messageChannel.sendToQueue(
                            queueName,
                            new Buffer(JSON.stringify(updateMessageMsg)),
                            { persistent: true }  // make sure messages are stored until ack'ed
                        );

                        res.set("Content-Type", "application/json");
                        res.json(newMessage);
                    })
                    .catch(err => {
                        console.log(err)
                        res.sendStatus(500).end()
                        return
                    })
            }
        })
        .catch(err => {
            console.log(err)
            res.sendStatus(500).end()
        })
});


// DELETE: If the current user isn't the creator of this message, respond with the status code 403 (Forbidden). Otherwise, delete the message and respond with a the plain text message indicating that the delete was successful.
app.delete("/v1/messages/:messageID", (req, res) => {
    let fetchedUser = req.get("X-User");
    let user = JSON.parse(fetchedUser);
    let messID = req.params.messageID

    Message.findById(messID).exec()

        .then(doc => {
            if (doc.creator.id != user.id) {
                res.sendStatus(403).end()
                return
            }
            Message.deleteOne({ _id: messID }, err => {
                if (err) {
                    console.log(err)
                    res.sendStatus(500).end()
                    return
                }
            });
            const updateMessageMsg = {
                "type": "message-delete",
                "data": messID,
                "userList": []
            };
            console.log(updateMessageMsg)
            messageChannel.sendToQueue(
                queueName,
                new Buffer(JSON.stringify(updateMessageMsg)),
                { persistent: true }  // make sure messages are stored until ack'ed
            );

            res.send("Delete was successful!")

        })
        .catch(err => {
            console.log(err)
            res.sendStatus(500).end()
        })
});

app.post("/v1/handleOffer", (req, res) => {

    const offerMsg = {
        "type": req.body.type,
        "data": req.body.data,
        "userList": req.body.userList
    };

    console.log(offerMsg)

    messageChannel.sendToQueue(
        queueName,
        new Buffer(JSON.stringify(offerMsg)),
        { persistent: true }  // make sure messages are stored until ack'ed
    );
    res.status(200).send("Offer Received")
})

setTimeout(runMicroServer, 20000);

function runMicroServer() {
amqp.connect('amqp://' + rabbitAddr, (err, conn) => {
    if (err) {
        console.log("Failed to connect to Rabbit Instance from API Server. due to ", err);
        process.exit(1);
    }

    mongoose.connect(url, { useNewUrlParser: true })
        .then(() => { console.log("Connected to database!") })
        .catch((err) => {
            console.log("couldnt connect to database because of ERR: ", err)
            process.exit(1);
        })

    conn.createChannel((err, ch) => {
        if (err) {
            console.log("Failed to create NewOrdersChannel from API Server");
            process.exit(1);
        }

        ch.assertQueue(queueName, { durable: true });
        messageChannel = ch;

        server.listen(port, () => {
            console.log(`Listening on port ${port}`);
        });
    });
});
}