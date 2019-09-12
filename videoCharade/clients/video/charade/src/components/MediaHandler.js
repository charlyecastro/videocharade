export default class MediaHandler {
    getPermissions(){
        return new Promise((res, rej) => {
            navigator.mediaDevices.getUserMedia({video: true, audio : true})
            .then((stream) => {
                res(stream)
            })
            .catch(err => {
                throw new Error(`Unable to fetch stream ${err}`)
                
            })
        })
    }
}


        // this.mediaHandler.getPermissions()
        //     .then((stream) => {
        //         this.setState({
        //             hasMedia: true,
        //             localStream: stream
        //         })
        //         try {
        //             this.myVideo.srcObject = stream
        //         } catch (e) {
        //             this.myVideo.src = URL.createObjectURL(stream)
        //         }
        //         this.myVideo.play()
        //     })