const mongoose = require("mongoose")

const channelSchema = mongoose.Schema({
    _id: mongoose.Schema.Types.ObjectId,
    name: { type: String, required: true },
    nameLower: String,
    description: String,
    private: { type: Boolean, default: false },
    members: [],
    createdAt: { type: Date, default: Date.now },
    creator: {},
    editedAt: { type: Date, default: Date.now },
})

module.exports = mongoose.model('Channel', channelSchema)