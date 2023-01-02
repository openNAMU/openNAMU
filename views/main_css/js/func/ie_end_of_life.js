function opennamu_do_ie_end_support() {
    if(document.currentScript === undefined) {
        window.location = 'microsoft-edge:' + window.location;
        setTimeout(function() {
            window.location = 'https://support.microsoft.com/office/160fa918-d581-4932-9e4e-1075c4713595';
        }, 1);
    }
}

opennamu_do_ie_end_support();
