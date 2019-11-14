function do_stop_exit() {
    window.onbeforeunload = function(){ return ''; }
}