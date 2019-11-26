function do_open_folding(data, element) { 
    var fol = document.getElementById(data);
    if(fol.style.display === '' || (fol.style.display === 'inline-block' || fol.style.display === 'block')) { 
        document.getElementById(data).style.display = 'none';
        element.innerHTML = '[+]'
    } else {
        document.getElementById(data).style.display = 'block';
        element.innerHTML = '[-]' 
    } 
}