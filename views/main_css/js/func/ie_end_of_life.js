function opennamu_do_ie_end_support() {
    if(document.currentScript === undefined) {
        window.location = 'microsoft-edge:' + window.location;
        setTimeout(function() {
            window.location = 'https://go.microsoft.com/fwlink/?linkid=2135547';
        }, 1);
    }
}

opennamu_do_ie_end_support();