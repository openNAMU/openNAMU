var save_data = '';
var open = 0;

function opening(data) {
    save_data = data;
    if(data === 'recent_cel') {
        var element = document.getElementById(data);
        var element_2 = document.getElementById('other_cel');
    } else {
        var element = document.getElementById(data);
        var element_2 = document.getElementById('recent_cel');
    }

    if(element.style.display == 'none') {
        element.style.display = 'block';
        element_2.style.display = 'none';
    } else {
        element.style.display = 'none';
    }

    open = 1;
    setTimeout(function() { open = 0; }, 100);
}

document.onclick = function(event) {
    for(var node = event.target; node != document.body; node = node.parentNode) {
        if(save_data !== '' && open == 0) {
            if(node.id === save_data) {
                break;
            } else {
                document.getElementById(save_data).style.display = 'none';
            }
        }
    }
}