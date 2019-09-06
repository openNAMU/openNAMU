function folding(num) {
  var fol = document.getElementById("folding_" + num);
  if (fol.style.display === "inline-block" || fol.style.display === "block") {
    fol.style.display = "none";
  } else {
    fol.style.display = "block";
  }
}
