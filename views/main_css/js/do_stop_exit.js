var go_save_zone = 0;

function do_stop_exit() {
    window.onbeforeunload = function() {
        console.log(go_save_zone);
        data = document.getElementById('content').value;
        origin = document.getElementById('origin').value;
        if(data !== origin && go_save_zone != 1) {
            return ''; 
        }
    }
}