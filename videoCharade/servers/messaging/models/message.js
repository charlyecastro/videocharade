const mongoose = require("mongoose")

const messageSchema = mongoose.Schema({
    _id: mongoose.Schema.Types.ObjectId,
    channelID: { type: mongoose.Schema.Types.ObjectId, ref: 'Channel', required: true },
    name: String,
    body: String,
    createdAt: { type: Date, default: Date.now },
    creator: {},
    editedAt: { type: Date, default: Date.now },
})

module.exports = mongoose.model('Message', messageSchema)