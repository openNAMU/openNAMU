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
    if(save_data !== '' && open == 0) {
        document.getElementById(save_data).style.display = 'none';
    }
}