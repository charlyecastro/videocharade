const http = require("http")

const authenticator = function (req, res, next) {
    let xUser = req.get('X-User');
    if (!xUser) {
        res.sendStatus(403).end()
        return
    } else {
        next()
    }
}

module.exports = authenticator