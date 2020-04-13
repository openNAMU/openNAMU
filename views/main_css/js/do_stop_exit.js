window.addEventListener('DOMContentLoaded', function() {
    if(window.location.pathname.match(/^\/edit\//i)) {
        window.onbeforeunload = function() {
            data = document.getElementById('content').value;
            origin = document.getElementById('origin').value;
            if(data !== origin) {
                return '';
            }
        }
    }
});

function save_stop_exit() {
    window.onbeforeunload = function () { }
}