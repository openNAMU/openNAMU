function do_open_folding(num, include_num, element) { 
    var fol = document.getElementById(include_num + 'folding_' + num); 
    if(fol.style.display === 'inline-block' || fol.style.display === 'block') { 
        fol.style.display = 'none';
        element.innerHTML = '[+]'
    } else {
        fol.style.display = 'block';
        element.innerHTML = '[-]' 
    } 
}