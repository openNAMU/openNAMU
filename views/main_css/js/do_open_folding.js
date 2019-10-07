function do_open_folding(num, include_num) { 
    var fol = document.getElementById(include_num + 'folding_' + num); 
    if(fol.style.display === 'inline-block' || fol.style.display === 'block') { 
        fol.style.display = 'none';
    } else {
        fol.style.display = 'block'; 
    } 
}