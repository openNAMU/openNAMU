function opening(data) {
    var element = document.getElementById(data);
    if(element.style.display == 'none') {
        element.style.display = 'block';
    } else {
        element.style.display = 'none';
    }
}