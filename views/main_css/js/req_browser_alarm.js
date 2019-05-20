// https://programmingsummaries.tistory.com/379
function req_alarm() {
    Notification.requestPermission(function (result) {
        if(result === 'denied') {
            document.cookie = 'topic_req=false;';
        } else {
            document.cookie = 'topic_req=false;';
        }
    });
}